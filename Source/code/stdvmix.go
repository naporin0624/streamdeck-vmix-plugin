package stdvmix

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/FlowingSPDG/streamdeck"
	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	vmixhttp "github.com/FlowingSPDG/vmix-go/http"
	"github.com/puzpuzpuz/xsync/v3"
)

const (
	// AppName Streamdeck plugin app name
	AppName = "dev.flowingspdg.vmix.sdPlugin"

	// ActionFunction SendFunction action Name
	ActionFunction = "dev.flowingspdg.vmix.function"

	// ActionPreview Preview input action Name
	ActionPreview = "dev.flowingspdg.vmix.preview"

	// ActionProgram Take input action Name
	ActionProgram = "dev.flowingspdg.vmix.program"
)

const (
	// tally color
	tallyInactive string = "data:image/png;base64,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TVA9Imh0dHA6Ly93d3cuZ2ltcC5vcmcveG1wLyIKICAgIHhtbG5zOmRjPSJodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyIKICAgIHhtbG5zOnhtcD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLyIKICAgeG1wTU06RG9jdW1lbnRJRD0iZ2ltcDpkb2NpZDpnaW1wOjViY2U0YWU3LTI5OTMtNDI0ZS04MDgwLWEzMzJjMTc2OGM4OCIKICAgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDo5Y2JiMTk3MS1mMmFiLTRlMDQtYjdmNy1hODAxZmRiMGE0NzMiCiAgIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpjZWM4Nzc0OC04MmVjLTRiOWYtOTg1MC1lNmJlNDY0MTJiZTYiCiAgIEdJTVA6QVBJPSIyLjAiCiAgIEdJTVA6UGxhdGZvcm09Ik1hYyBPUyIKICAgR0lNUDpUaW1lU3RhbXA9IjE2MTk2NjUxMTA5ODcyNjQiCiAgIEdJTVA6VmVyc2lvbj0iMi4xMC4xNCIKICAgZGM6Rm9ybWF0PSJpbWFnZS9wbmciCiAgIHhtcDpDcmVhdG9yVG9vbD0iR0lNUCAyLjEwIj4KICAgPGlwdGNFeHQ6TG9jYXRpb25DcmVhdGVkPgogICAgPHJkZjpCYWcvPgogICA8L2lwdGNFeHQ6TG9jYXRpb25DcmVhdGVkPgogICA8aXB0Y0V4dDpMb2NhdGlvblNob3duPgogICAgPHJkZjpCYWcvPgogICA8L2lwdGNFeHQ6TG9jYXRpb25TaG93bj4KICAgPGlwdGNFeHQ6QXJ0d29ya09yT2JqZWN0PgogICAgPHJkZjpCYWcvPgogICA8L2lwdGNFeHQ6QXJ0d29ya09yT2JqZWN0PgogICA8aXB0Y0V4dDpSZWdpc3RyeUlkPgogICAgPHJkZjpCYWcvPgogICA8L2lwdGNFeHQ6UmVnaXN0cnlJZD4KICAgPHhtcE1NOkhpc3Rvcnk+CiAgICA8cmRmOlNlcT4KICAgICA8cmRmOmxpCiAgICAgIHN0RXZ0OmFjdGlvbj0ic2F2ZWQiCiAgICAgIHN0RXZ0OmNoYW5nZWQ9Ii8iCiAgICAgIHN0RXZ0Omluc3RhbmNlSUQ9InhtcC5paWQ6MDNhZmM1ZDMtZGI4ZC00NjA4LTliN2UtNDQwNzFmMzY3YWUxIgogICAgICBzdEV2dDpzb2Z0d2FyZUFnZW50PSJHaW1wIDIuMTAgKE1hYyBPUykiCiAgICAgIHN0RXZ0OndoZW49IjIwMjEtMDQtMjlUMTE6NTg6MzArMDk6MDAiLz4KICAgIDwvcmRmOlNlcT4KICAgPC94bXBNTTpIaXN0b3J5PgogICA8cGx1czpJbWFnZVN1cHBsaWVyPgogICAgPHJkZjpTZXEvPgogICA8L3BsdXM6SW1hZ2VTdXBwbGllcj4KICAgPHBsdXM6SW1hZ2VDcmVhdG9yPgogICAgPHJkZjpTZXEvPgogICA8L3BsdXM6SW1hZ2VDcmVhdG9yPgogICA8cGx1czpDb3B5cmlnaHRPd25lcj4KICAgIDxyZGY6U2VxLz4KICAgPC9wbHVzOkNvcHlyaWdodE93bmVyPgogICA8cGx1czpMaWNlbnNvcj4KICAgIDxyZGY6U2VxLz4KICAgPC9wbHVzOkxpY2Vuc29yPgogIDwvcmRmOkRlc2NyaXB0aW9uPgogPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgIAo8P3hwYWNrZXQgZW5kPSJ3Ij8+7MRfwQAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB+UEHQI6HgGMPmcAAABySURBVHja7dAxEQAwCASwUuUY/zsUsDMkElJJ+rH6CgQJEiRIkCBBghAkSJAgQYIECUKQIEGCBAkSJAhBggQJEiRIkCBBCBIkSJAgQYIEIUiQIEGCBAkShCBBggQJEiRIkCAECRIkSJAgQYIQJEiQoCsG1+IEBwGJzGQAAAAASUVORK5CYII="
	tallyPreview  string = "data:image/png;base64,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ICB4bWxuczpHSU1QPSJodHRwOi8vd3d3LmdpbXAub3JnL3htcC8iCiAgICB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iCiAgICB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iCiAgIHhtcE1NOkRvY3VtZW50SUQ9ImdpbXA6ZG9jaWQ6Z2ltcDo3ZjI1NjAzOC0xOTkwLTQ5Y2MtOTVlMi1jNzI3NDBjYzYxNTAiCiAgIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6NDk4N2VkNDktMTZhMC00NjA4LWE4NzItYzNiN2ZmMmY0ZDlhIgogICB4bXBNTTpPcmlnaW5hbERvY3VtZW50SUQ9InhtcC5kaWQ6YjIwOWVhZGMtZjFmNC00MDQxLWE5NzMtYTkwZGFhYTdhOTUzIgogICBHSU1QOkFQST0iMi4wIgogICBHSU1QOlBsYXRmb3JtPSJNYWMgT1MiCiAgIEdJTVA6VGltZVN0YW1wPSIxNjE5NjY1MDQxMDUyOTAxIgogICBHSU1QOlZlcnNpb249IjIuMTAuMTQiCiAgIGRjOkZvcm1hdD0iaW1hZ2UvcG5nIgogICB4bXA6Q3JlYXRvclRvb2w9IkdJTVAgMi4xMCI+CiAgIDxpcHRjRXh0OkxvY2F0aW9uQ3JlYXRlZD4KICAgIDxyZGY6QmFnLz4KICAgPC9pcHRjRXh0OkxvY2F0aW9uQ3JlYXRlZD4KICAgPGlwdGNFeHQ6TG9jYXRpb25TaG93bj4KICAgIDxyZGY6QmFnLz4KICAgPC9pcHRjRXh0OkxvY2F0aW9uU2hvd24+CiAgIDxpcHRjRXh0OkFydHdvcmtPck9iamVjdD4KICAgIDxyZGY6QmFnLz4KICAgPC9pcHRjRXh0OkFydHdvcmtPck9iamVjdD4KICAgPGlwdGNFeHQ6UmVnaXN0cnlJZD4KICAgIDxyZGY6QmFnLz4KICAgPC9pcHRjRXh0OlJlZ2lzdHJ5SWQ+CiAgIDx4bXBNTTpIaXN0b3J5PgogICAgPHJkZjpTZXE+CiAgICAgPHJkZjpsaQogICAgICBzdEV2dDphY3Rpb249InNhdmVkIgogICAgICBzdEV2dDpjaGFuZ2VkPSIvIgogICAgICBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOjdjNzkzZjA3LTViNjQtNDc0ZS04Mjk3LWYzMTFlOTczMDkwYyIKICAgICAgc3RFdnQ6c29mdHdhcmVBZ2VudD0iR2ltcCAyLjEwIChNYWMgT1MpIgogICAgICBzdEV2dDp3aGVuPSIyMDIxLTA0LTI5VDExOjU3OjIxKzA5OjAwIi8+CiAgICA8L3JkZjpTZXE+CiAgIDwveG1wTU06SGlzdG9yeT4KICAgPHBsdXM6SW1hZ2VTdXBwbGllcj4KICAgIDxyZGY6U2VxLz4KICAgPC9wbHVzOkltYWdlU3VwcGxpZXI+CiAgIDxwbHVzOkltYWdlQ3JlYXRvcj4KICAgIDxyZGY6U2VxLz4KICAgPC9wbHVzOkltYWdlQ3JlYXRvcj4KICAgPHBsdXM6Q29weXJpZ2h0T3duZXI+CiAgICA8cmRmOlNlcS8+CiAgIDwvcGx1czpDb3B5cmlnaHRPd25lcj4KICAgPHBsdXM6TGljZW5zb3I+CiAgICA8cmRmOlNlcS8+CiAgIDwvcGx1czpMaWNlbnNvcj4KICA8L3JkZjpEZXNjcmlwdGlvbj4KIDwvcmRmOlJERj4KPC94OnhtcG1ldGE+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAKPD94cGFja2V0IGVuZD0idyI/Plfb3jAAAAAGYktHRAD/AP8A/6C9p5MAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQflBB0CORW9c7QsAAAAcElEQVR42u3QMQEAAAgDoGn/zprA3wMiUJlMOLUCQYIECRIkSJAgBAkSJEiQIEGCECRIkCBBggQJQpAgQYIECRIkSBCCBAkSJEiQIEEIEiRIkCBBggQhSJAgQYIECRIkCEGCBAkSJEiQIAQJEiToiwUf1QKOQh77lQAAAABJRU5ErkJggg=="
	tallyProgram  string = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEgAAABICAYAAABV7bNHAAABhGlDQ1BJQ0MgcHJvZmlsZQAAKJF9kT1Iw0AcxV/TSkUqHewgopChOlkQFemoVShChVArtOpgcukXNGlIUlwcBdeCgx+LVQcXZ10dXAVB8APEzc1J0UVK/F9aaBHjwXE/3t173L0DhEaFaVZgAtB020wnE2I2tyoGXxHACATEEZaZZcxJUgqe4+sePr7exXiW97k/R7+atxjgE4lnmWHaxBvEM5u2wXmfOMJKskp8Tjxu0gWJH7mutPiNc9FlgWdGzEx6njhCLBa7WOliVjI14mniqKrplC9kW6xy3uKsVWqsfU/+wlBeX1nmOs1hJLGIJUgQoaCGMiqwEaNVJ8VCmvYTHv4h1y+RSyFXGYwcC6hCg+z6wf/gd7dWYWqylRRKAD0vjvMxCgR3gWbdcb6PHad5AvifgSu94682gPgn6fWOFj0CwtvAxXVHU/aAyx1g8MmQTdmV/DSFQgF4P6NvygEDt0DfWqu39j5OH4AMdZW6AQ4OgbEiZa97vLu3u7d/z7T7+wF1rnKoxhB+yAAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB+UEHQI4IYXccdgAAABwSURBVHja7dAxAQAACAOgaf/OmsDfAyJQk0w4tQJBggQJEiRIkCAECRIkSJAgQYIQJEiQIEGCBAlCkCBBggQJEiRIEIIECRIkSJAgQQgSJEiQIEGCBCFIkCBBggQJEiQIQYIECRIkSJAgBAkSJOiLBSDUAo5LcSa/AAAAAElFTkSuQmCC"
)

type Input struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Number int    `json:"number"`
}

type vMix struct {
	client *vmixhttp.Client
	inputs []Input
}

type StdVmix struct {
	// logger
	logger *log.Logger
	// StreamDeck Client
	c *streamdeck.Client

	// vMix Clients
	// TODO: 削除/設定が変更されたときにK/Vからも削除する
	vMixClients *vMixConnections

	// Contexts
	sendFuncContexts *xsync.MapOf[string, SendFunctionPI]
	previewContexts  *xsync.MapOf[string, PreviewPI]
	programContexts  *xsync.MapOf[string, ProgramPI]
}

func NewStdVmix(ctx context.Context, params streamdeck.RegistrationParams, logWriter io.Writer) *StdVmix {
	logger := log.New(os.Stdout, "vMix[FlowingSPDG]: ", log.LstdFlags)
	logger.SetOutput(io.MultiWriter(logWriter, os.Stdout))
	logger.SetFlags(log.Ldate | log.Ltime)

	logger.Println("Initiating new vMix plugin instance...")

	client := streamdeck.NewClient(ctx, params)
	ret := &StdVmix{
		logger:           logger,
		c:                client,
		vMixClients:      newVMixConnections(),
		sendFuncContexts: xsync.NewMapOf[string, SendFunctionPI](),
		previewContexts:  xsync.NewMapOf[string, PreviewPI](),
		programContexts:  xsync.NewMapOf[string, ProgramPI](),
	}

	actionFunc := client.Action(ActionFunction)
	actionFunc.RegisterHandler(streamdeck.WillAppear, ret.SendFuncWillAppearHandler)
	actionFunc.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.sendFuncContexts.Delete(event.Context)
		return nil
	})
	actionFunc.RegisterHandler(streamdeck.KeyDown, ret.SendFuncKeyDownHandler)
	actionFunc.RegisterHandler(streamdeck.DidReceiveSettings, ret.SendFuncDidReceiveSettingsHandler)

	actionPrev := client.Action(ActionPreview)
	actionPrev.RegisterHandler(streamdeck.WillAppear, ret.PreviewWillAppearHandler)
	actionPrev.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.previewContexts.Delete(event.Context)
		return nil
	})
	actionPrev.RegisterHandler(streamdeck.KeyDown, ret.PreviewKeyDownHandler)
	actionPrev.RegisterHandler(streamdeck.DidReceiveSettings, ret.PreviewDidReceiveSettingsHandler)

	actionProgram := client.Action(ActionProgram)
	actionProgram.RegisterHandler(streamdeck.WillAppear, ret.ProgramWillAppearHandler)
	actionProgram.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.programContexts.Delete(event.Context)
		return nil
	})
	actionProgram.RegisterHandler(streamdeck.KeyDown, ret.ProgramKeyDownHandler)
	actionProgram.RegisterHandler(streamdeck.DidReceiveSettings, ret.ProgramDidReceiveSettingsHandler)

	ret.c = client

	return ret
}

type InputsForPI struct {
	Inputs []Input `json:"inputs"`
}

type SendToPropertyInspectorPayload[T any] struct {
	Event   string `json:"event"`
	Payload T      `json:"payload"`
}

// Update inputs Contextの数だけ更新が入るので負荷が高いかもしれない
func (s *StdVmix) Update() {
	// now := time.Now()
	// s.logger.Println("Updating")

	// vMixの更新
	activeKeys := make([]vMixKey, 0, s.previewContexts.Size()+s.programContexts.Size())
	s.previewContexts.Range(func(ctxStr string, pi PreviewPI) bool {
		activeKeys = append(activeKeys, vMixKey{host: pi.Host, port: pi.Port})
		return true
	})
	s.programContexts.Range(func(ctxStr string, pi ProgramPI) bool {
		activeKeys = append(activeKeys, vMixKey{host: pi.Host, port: pi.Port})
		return true
	})
	s.vMixClients.UpdateVMixes(activeKeys)

	// PRVの更新
	// s.logger.Printf("Updating %d PRV contexts\n", s.previewContexts.Size())
	s.previewContexts.Range(func(ctxStr string, pi PreviewPI) bool {
		ctx := context.Background()
		ctx = sdcontext.WithContext(ctx, ctxStr)

		go func() {
			// inputの更新
			v, err := s.vMixClients.loadOrStore(pi.Host, pi.Port)
			if err != nil {
				return
			}
			s.c.SendToPropertyInspector(ctx, SendToPropertyInspectorPayload[InputsForPI]{
				Event: "inputs",
				Payload: InputsForPI{
					Inputs: v.inputs,
				},
			})

			// TALLYの更新、不要なら飛ばす
			// TODO: 関数に分ける
			if !pi.Tally {
				return
			}
			currentPreview := v.client.Preview
			// Mixの場合
			if pi.Mix > 1 {
				for _, mix := range v.client.Mix {
					if int(mix.Number) == pi.Mix {
						currentPreview = mix.Preview
						break
					}
				}
			}
			// TODO: 毎回SetImageをしたくないので、状態管理して変更時のみトリガーする
			tally := false
			for _, i := range v.inputs {
				if i.Key == pi.Input && currentPreview == uint(i.Number) {
					tally = true

					break
				}
			}
			if tally {
				if err := s.c.SetImage(ctx, tallyPreview, streamdeck.HardwareAndSoftware); err != nil {
					s.c.LogMessage(fmt.Sprintf("failed to set preview tally: %v", err))
				}
			} else {
				if err := s.c.SetImage(ctx, tallyInactive, streamdeck.HardwareAndSoftware); err != nil {
					s.c.LogMessage(fmt.Sprintf("failed to set preview tally: %v", err))
				}
			}
		}()
		return true
	})

	// PGMの更新
	// s.logger.Printf("Updating %d PGM contexts\n", s.programContexts.Size())
	s.programContexts.Range(func(ctxStr string, pi ProgramPI) bool {
		ctx := context.Background()
		ctx = sdcontext.WithContext(ctx, ctxStr)

		go func() {
			v, err := s.vMixClients.loadOrStore(pi.Host, pi.Port)
			if err != nil {
				return
			}
			s.c.SendToPropertyInspector(ctx, SendToPropertyInspectorPayload[InputsForPI]{
				Event: "inputs",
				Payload: InputsForPI{
					Inputs: v.inputs,
				},
			})

			// TALLYの更新、不要なら飛ばす
			// TODO: 関数に分ける
			if !pi.Tally {
				return
			}
			activeInput := v.client.Active
			// Mixの場合
			if pi.Mix > 1 {
				for _, mix := range v.client.Mix {
					if int(mix.Number) == pi.Mix {
						activeInput = mix.Preview
						break
					}
				}
			}
			// TODO: 毎回SetImageをしたくないので、状態管理して変更時のみトリガーする
			tally := false
			for _, i := range v.inputs {
				if i.Key == pi.Input && activeInput == uint(i.Number) {
					tally = true
					break
				}
			}
			if tally {
				if err := s.c.SetImage(ctx, tallyProgram, streamdeck.HardwareAndSoftware); err != nil {
					s.c.LogMessage(fmt.Sprintf("failed to set preview tally: %v", err))
				}
			} else {
				if err := s.c.SetImage(ctx, tallyInactive, streamdeck.HardwareAndSoftware); err != nil {
					s.c.LogMessage(fmt.Sprintf("failed to set preview tally: %v", err))
				}
			}
		}()
		return true
	})

	// s.logger.Printf("Updated in %v\n", time.Since(now))
}

func (s *StdVmix) Run(ctx context.Context) error {
	go func() {
		for {
			time.Sleep(time.Second / 5) // 0.2s
			select {
			case <-ctx.Done():
				return
			default:
				s.Update()
			}
		}
	}()
	return s.c.Run()
}
