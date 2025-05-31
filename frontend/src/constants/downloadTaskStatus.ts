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
}

