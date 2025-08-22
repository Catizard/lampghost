export enum DownloadTaskStatus {
  PREPARE = "Prepare",
  DOWNLOAD = "Download",
  SUCCESS = "Success",
  ERROR = "Error",
  CANCEL = "Cancel"
}

export namespace DownloadTaskStatus {
  export function from(v: number): DownloadTaskStatus {
    switch (v) {
      case 0: return DownloadTaskStatus.PREPARE;
      case 1: return DownloadTaskStatus.DOWNLOAD;
      case 2: return DownloadTaskStatus.SUCCESS;
      case 3: return DownloadTaskStatus.ERROR;
      case 4: return DownloadTaskStatus.CANCEL;
      default: throw "unexpected DownloadTaskStatus: " + v;
    }
  }

  export function restartable(v: number | DownloadTaskStatus): boolean {
    const status: DownloadTaskStatus = typeof (v) == "number" ? DownloadTaskStatus.from(v) : v;
    return status == DownloadTaskStatus.CANCEL || status == DownloadTaskStatus.ERROR;
  }

  export function cancelable(v: number | DownloadTaskStatus): boolean {
    const status: DownloadTaskStatus = typeof (v) == "number" ? DownloadTaskStatus.from(v) : v;
    return status == DownloadTaskStatus.DOWNLOAD;
  }

  export function compare(a: DownloadTaskStatus, b: DownloadTaskStatus): number {
    const to = (s: DownloadTaskStatus) => {
      switch (s) {
        case DownloadTaskStatus.DOWNLOAD: return 0;
        case DownloadTaskStatus.ERROR: return 1;
        case DownloadTaskStatus.PREPARE: return 2;
        case DownloadTaskStatus.SUCCESS: return 3;
        case DownloadTaskStatus.CANCEL: return 4;
      }
    }
    return to(a) - to(b);
  }
}

