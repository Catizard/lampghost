export enum DownloadTaskStatus {
  PREPARE = "Prepare",
  DOWNLOAD = "Download",
  SUCCESS = "Success",
  ERROR = "Error"
}

export namespace DownloadTaskStatus {
  export function from(v: number): DownloadTaskStatus {
    switch (v) {
      case 0: return DownloadTaskStatus.PREPARE;
      case 1: return DownloadTaskStatus.DOWNLOAD;
      case 2: return DownloadTaskStatus.SUCCESS;
      case 3: return DownloadTaskStatus.ERROR;
    }
  }
}

