import { DialogApiInjection } from "naive-ui/es/dialog/src/DialogProvider";
import { LoadingBarApiInjection } from "naive-ui/es/loading-bar/src/LoadingBarProvider";
import { MessageApiInjection } from "naive-ui/es/message/src/MessageProvider";
import { ModalApiInjection } from "naive-ui/es/modal/src/ModalProvider";
import { NotificationApiInjection } from "naive-ui/es/notification/src/NotificationProvider";

export { };

declare global {
  interface Window {
    $message: MessageApiInjection,
    $dialog: DialogApiInjection,
    $notification: NotificationApiInjection,
    $loadingBar: LoadingBarApiInjection,
    $modal: ModalApiInjection,
    $notifySuccess(content: string): void;
    $notifyDefault(content: string): void;
    $notifyInfo(content: string): void;
    $notifyWarning(content: string): void;
    $notifyError(content: string | Error): void;
  }
}
