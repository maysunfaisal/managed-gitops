package eventloop

import (
	"github.com/maysunfaisal/managed-gitops/backend/apis/managed-gitops/v1alpha1"
	"github.com/maysunfaisal/managed-gitops/backend/eventloop/eventlooptypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Controller Event Loop Test", func() {

	Context("Controller event loop responds to channel events", func() {

		It("Should pass events received on input channel to mock workspace event loop router", func() {

			mockOutputChannelFactory := &mockWorkspaceEventLoopFactory{
				mockChannel: make(chan eventlooptypes.EventLoopMessage),
			}

			loop := newControllerEventLoopWithFactory(mockOutputChannelFactory)

			loop.eventLoopInputChannel <- eventlooptypes.EventLoopEvent{
				EventType: eventlooptypes.DeploymentModified,
				Request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "",
						Name:      "",
					},
				},
				Client:                  nil,
				ReqResource:             v1alpha1.GitOpsDeploymentTypeName,
				AssociatedGitopsDeplUID: "",
				WorkspaceID:             "",
			}

			resp := <-mockOutputChannelFactory.mockChannel
			Expect(resp).NotTo(BeNil())

		})
	})
})

// mockWorkspaceEventLoopFactory is a mock of workspaceEventLoopRouterFactory
type mockWorkspaceEventLoopFactory struct {
	mockChannel chan eventlooptypes.EventLoopMessage
}

var _ workspaceEventLoopRouterFactory = &mockWorkspaceEventLoopFactory{}

func (cetf *mockWorkspaceEventLoopFactory) startWorkspaceEventLoopRouter(workspaceID string) chan eventlooptypes.EventLoopMessage {
	// Rather than starting a new workspace event loop, instead just return a pre-provided channel
	return cetf.mockChannel
}
