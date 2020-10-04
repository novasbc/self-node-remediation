package controllers

import (
	"context"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	machinev1beta1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"time"
)

var _ = Describe("Machine Controller", func() {
	Context("Unhealthy machine with api-server access", func() {
		logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

		machineName := "machine1"
		machine1 := &machinev1beta1.Machine{}
		machineNamespacedName := types.NamespacedName{
			Name:      machineName,
			Namespace: machineNamespace,
		}

		It("Check the machine exists", func() {
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), machineNamespacedName, machine1)
			}, 10*time.Second, 250*time.Millisecond).Should(BeNil())

			Expect(machine1.Name).To(Equal("machine1"))
			Expect(machine1.Status.NodeRef).ToNot(BeNil())
		})

		It("Mark machine as unhealthy", func() {
			if machine1.Annotations == nil {
				machine1.Annotations = make(map[string]string)
			}
			machine1.Annotations[externalRemediationAnnotation] = ""
			err := k8sClient.Update(context.TODO(), machine1)
			Expect(err).ToNot(HaveOccurred())
		})

		node := &v1.Node{}
		nodeNamespacedName := client.ObjectKey{
			Name:      "node1",
			Namespace: "",
		}

		It("Verify that node was marked as unschedulable ", func() {
			Eventually(func() bool {
				node = &v1.Node{}
				Expect(k8sClient.Get(context.TODO(), nodeNamespacedName, node)).To(Succeed())

				return node.Spec.Unschedulable

			}, 5*time.Second, 250*time.Millisecond).Should(BeTrue())
		})

		It("Add unshedulable taint to node to simulate node controller", func() {
			node.Spec.Taints = append(node.Spec.Taints, *NodeUnschedulableTaint)
			Expect(k8sClient.Update(context.TODO(), node)).To(Succeed())
		})

		It("Verify that time has been added to annotation", func() {
			Eventually(func() string {
				machine1 = &machinev1beta1.Machine{}
				Expect(k8sClient.Get(context.TODO(), machineNamespacedName, machine1)).To(Succeed())

				//give some time to the machine controller to update the time in the annotation
				return machine1.Annotations[externalRemediationAnnotation]

			}, 5*time.Second, 250*time.Millisecond).ShouldNot(BeEmpty())
		})

		It("Verify that node backup annotation matches the node", func() {
			Expect(machine1.Annotations).To(HaveKey(nodeBackupAnnotation))
			Expect(machine1.Annotations[nodeBackupAnnotation]).ToNot(BeEmpty())
			nodeToRestore := &v1.Node{}
			Expect(json.Unmarshal([]byte(machine1.Annotations[nodeBackupAnnotation]), nodeToRestore)).To(Succeed())

			node = &v1.Node{}
			Expect(k8sClient.Get(context.TODO(), nodeNamespacedName, node)).To(Succeed())

			//todo why do we need the following 2 lines? this might be a bug
			nodeToRestore.TypeMeta.Kind = "Node"
			nodeToRestore.TypeMeta.APIVersion = "v1"
			Expect(nodeToRestore).To(Equal(node))
		})

		It("Verify that watchdog is not receiving food", func() {
			currentLastFoodTime := dummyDog.GetLastFoodTime()
			Consistently(func() time.Time {
				return dummyDog.GetLastFoodTime()
			}, 5*reconcileInterval, 1*time.Second).Should(Equal(currentLastFoodTime))
		})

		now := time.Now()
		It("Update annotation time to accelerate the progress", func() {
			oldTime := now.Add(-safeTimeToAssumeNodeRebooted).Add(-time.Minute)
			machine1.Annotations[externalRemediationAnnotation] = oldTime.Format(time.RFC3339)
			Expect(k8sClient.Update(context.TODO(), machine1)).To(Succeed())
		})

		It("Verify that node has been deleted", func() {
			// in real world scenario, other machines will take care for the rest of the test but
			// in this test, we trick the machine to recover itself after we already verified it
			// tried to reboot
			shouldReboot = false

			Eventually(func() metav1.StatusReason {
				node = &v1.Node{}
				err := k8sClient.Get(context.TODO(), nodeNamespacedName, node)
				return errors.ReasonForError(err)
			}, 2*time.Second, 20*time.Millisecond).Should(Equal(metav1.StatusReasonNotFound))
		})

		It("Verify that node has been restored", func() {
			node = &v1.Node{}

			Eventually(func() error {
				return k8sClient.Get(context.TODO(), nodeNamespacedName, node)
			}, 5*time.Second, 250*time.Millisecond).Should(BeNil())

			Expect(node.CreationTimestamp.After(now)).To(BeTrue())
		})

		It("Verify that node is not marked as unschedulable", func() {
			Expect(node.Spec.Unschedulable).To(BeFalse())
		})

		It("Verify unhealthy annotation was removed", func() {
			Eventually(func() map[string]string {
				machine1 = &machinev1beta1.Machine{}
				Expect(k8sClient.Get(context.TODO(), machineNamespacedName, machine1)).To(Succeed())
				return machine1.Annotations
			}, 5*time.Second, 250*time.Millisecond).ShouldNot(HaveKey(externalRemediationAnnotation))

		})
	})
})
