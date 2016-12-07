package cmd_test

import (
	. "github.com/sclevine/cflocal/cf/cmd"
	"github.com/sclevine/cflocal/cf/cmd/mocks"
	"github.com/sclevine/cflocal/local"
	sharedmocks "github.com/sclevine/cflocal/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Export", func() {
	var (
		mockCtrl   *gomock.Controller
		mockUI     *sharedmocks.MockUI
		mockStager *mocks.MockStager
		mockRunner *mocks.MockRunner
		mockFS     *mocks.MockFS
		mockHelp   *mocks.MockHelp
		mockConfig *mocks.MockConfig
		cmd        *Export
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUI = sharedmocks.NewMockUI()
		mockStager = mocks.NewMockStager(mockCtrl)
		mockRunner = mocks.NewMockRunner(mockCtrl)
		mockFS = mocks.NewMockFS(mockCtrl)
		mockHelp = mocks.NewMockHelp(mockCtrl)
		mockConfig = mocks.NewMockConfig(mockCtrl)
		cmd = &Export{
			UI:     mockUI,
			Stager: mockStager,
			Runner: mockRunner,
			FS:     mockFS,
			Help:   mockHelp,
			Config: mockConfig,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Match", func() {
		It("should return true when the first argument is export", func() {
			Expect(cmd.Match([]string{"export"})).To(BeTrue())
			Expect(cmd.Match([]string{"not-export"})).To(BeFalse())
			Expect(cmd.Match([]string{})).To(BeFalse())
			Expect(cmd.Match(nil)).To(BeFalse())
		})
	})

	Describe("#Run", func() {
		It("should export a droplet as a Docker image", func() {
			droplet := newMockBufferCloser(mockCtrl)
			launcher := newMockBufferCloser(mockCtrl)
			localYML := &local.LocalYML{
				Applications: []*local.AppConfig{
					{Name: "some-other-app"},
					{
						Name: "some-app",
						Env:  map[string]string{"a": "b"},
					},
				},
			}
			gomock.InOrder(
				mockFS.EXPECT().ReadFile("./some-app.droplet").Return(droplet, int64(100), nil),
				mockStager.EXPECT().Launcher().Return(launcher, int64(200), nil),
				mockConfig.EXPECT().Load().Return(localYML, nil),
				mockRunner.EXPECT().Export(&local.RunConfig{
					Droplet:      droplet,
					DropletSize:  int64(100),
					Launcher:     launcher,
					LauncherSize: int64(200),
					AppConfig: &local.AppConfig{
						Name: "some-app",
						Env:  map[string]string{"a": "b"},
					},
				}, "some-reference").Return("some-id", nil),
				launcher.EXPECT().Close(),
				droplet.EXPECT().Close(),
			)
			Expect(cmd.Run([]string{"export", "-r", "some-reference", "some-app"})).To(Succeed())
			Expect(mockUI.Out).To(gbytes.Say("Exported some-app as some-reference with ID: some-id"))
		})

		// test without reference
	})
})
