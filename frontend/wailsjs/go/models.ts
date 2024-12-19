export namespace dto {
	
	export class DiffTableHeaderDto {
	    ID: number;
	    HeaderUrl: string;
	    DataUrl: string;
	    Name: string;
	    OriginalUrl?: string;
	    Symbol: string;
	    Contents: entity.DiffTableData[];
	
	    static createFrom(source: any = {}) {
	        return new DiffTableHeaderDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.HeaderUrl = source["HeaderUrl"];
	        this.DataUrl = source["DataUrl"];
	        this.Name = source["Name"];
	        this.OriginalUrl = source["OriginalUrl"];
	        this.Symbol = source["Symbol"];
	        this.Contents = this.convertValues(source["Contents"], entity.DiffTableData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace entity {
	
	export class DiffTableData {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    HeaderID: number;
	    Artist: string;
	    Comment: string;
	    Level: string;
	    lr2_bmdid: string;
	    Md5: string;
	    NameDiff: string;
	    Title: string;
	    url: string;
	    url_diff: string;
	    Sha256: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffTableData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.HeaderID = source["HeaderID"];
	        this.Artist = source["Artist"];
	        this.Comment = source["Comment"];
	        this.Level = source["Level"];
	        this.lr2_bmdid = source["lr2_bmdid"];
	        this.Md5 = source["Md5"];
	        this.NameDiff = source["NameDiff"];
	        this.Title = source["Title"];
	        this.url = source["url"];
	        this.url_diff = source["url_diff"];
	        this.Sha256 = source["Sha256"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DiffTableHeader {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    HeaderUrl: string;
	    data_url: string;
	    name: string;
	    original_url?: string;
	    symbol: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffTableHeader(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.HeaderUrl = source["HeaderUrl"];
	        this.data_url = source["data_url"];
	        this.name = source["name"];
	        this.original_url = source["original_url"];
	        this.symbol = source["symbol"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RivalInfo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    Name: string;
	    ScoreLogPath?: string;
	    SongDataPath?: string;
	    PlayCount: number;
	    MainUser: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RivalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.Name = source["Name"];
	        this.ScoreLogPath = source["ScoreLogPath"];
	        this.SongDataPath = source["SongDataPath"];
	        this.PlayCount = source["PlayCount"];
	        this.MainUser = source["MainUser"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace result {
	
	export class RtnData {
	    Data: any;
	    Code: number;
	    Msg: string;
	    // Go type: time
	    Timestamp: any;
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new RtnData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = source["Data"];
	        this.Code = source["Code"];
	        this.Msg = source["Msg"];
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	        this.Err = source["Err"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RtnDataList {
	    Rows: any[];
	    Code: number;
	    Msg: string;
	    // Go type: time
	    Timestamp: any;
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new RtnDataList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Rows = source["Rows"];
	        this.Code = source["Code"];
	        this.Msg = source["Msg"];
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	        this.Err = source["Err"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RtnMessage {
	    Code: number;
	    Msg: string;
	    // Go type: time
	    Timestamp: any;
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new RtnMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Code = source["Code"];
	        this.Msg = source["Msg"];
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	        this.Err = source["Err"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

