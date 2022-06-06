package application_event_loop

import (
	"github.com/maysunfaisal/managed-gitops/backend/apis/managed-gitops/v1alpha1"
	"github.com/maysunfaisal/managed-gitops/backend/eventloop/eventlooptypes"
	"github.com/maysunfaisal/managed-gitops/backend/eventloop/shared_resource_loop"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("ApplicationEventLoop Test", func() {

	Context("Application Event Loop test", func() {

		It("should forward messages received on the input channel, to the runner", func() {

			mockApplicationEventLoopRunnerFactory := mockApplicationEventLoopRunnerFactory{
				mockChannel: make(chan *eventlooptypes.EventLoopEvent),
			}

			inputChan := startApplicationEventQueueLoopWithFactory("", "", nil, &mockApplicationEventLoopRunnerFactory)

			inputChan <- eventlooptypes.EventLoopMessage{
				MessageType: eventlooptypes.ApplicationEventLoopMessageType_Event,
				Event: &eventlooptypes.EventLoopEvent{
					EventType:               eventlooptypes.DeploymentModified,
					Request:                 reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: ""}},
					Client:                  nil,
					ReqResource:             v1alpha1.GitOpsDeploymentTypeName,
					AssociatedGitopsDeplUID: "",
					WorkspaceID:             "",
				},
				ShutdownSignalled: false,
			}

			outputEvent := <-mockApplicationEventLoopRunnerFactory.mockChannel
			Expect(outputEvent).NotTo(BeNil())
			Expect(outputEvent.EventType).To(Equal(outputEvent.EventType))
			Expect(outputEvent.Request).To(Equal(outputEvent.Request))
			Expect(outputEvent.ReqResource).To(Equal(outputEvent.ReqResource))
			Expect(outputEvent.AssociatedGitopsDeplUID).To(Equal(outputEvent.AssociatedGitopsDeplUID))
			Expect(outputEvent.WorkspaceID).To(Equal(outputEvent.WorkspaceID))

		})
	})

})

// mockApplicationEventLoopRunnerFactory which returns a pre-provided channel, rather than starting a new goroutine.
type mockApplicationEventLoopRunnerFactory struct {
	mockChannel chan *eventlooptypes.EventLoopEvent
}

var _ applicationEventRunnerFactory = &mockApplicationEventLoopRunnerFactory{}

func (fact *mockApplicationEventLoopRunnerFactory) createNewApplicationEventLoopRunner(informWorkCompleteChan chan eventlooptypes.EventLoopMessage,
	sharedResourceEventLoop *shared_resource_loop.SharedResourceEventLoop,
	gitopsDeplUID string, workspaceID string, debugContext string) chan *eventlooptypes.EventLoopEvent {

	return fact.mockChannel

}
