export namespace config {
	
	export class ApplicationConfig {
	    UserName: string;
	    ScorelogFilePath: string;
	    SongdataFilePath: string;
	    ScoreFilePath: string;
	    InternalServerPort: number;
	    FolderSymbol: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplicationConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UserName = source["UserName"];
	        this.ScorelogFilePath = source["ScorelogFilePath"];
	        this.SongdataFilePath = source["SongdataFilePath"];
	        this.ScoreFilePath = source["ScoreFilePath"];
	        this.InternalServerPort = source["InternalServerPort"];
	        this.FolderSymbol = source["FolderSymbol"];
	    }
	}

}

export namespace dto {
	
	export class CourseInfoDto {
	    ID: number;
	    HeaderID: number;
	    Name: string;
	    Md5: string[];
	    Md5s: string;
	    Sha256: string[];
	    Sha256s: string;
	    NoSepJoinedSha256s: string;
	    Constranints: string;
	
	    static createFrom(source: any = {}) {
	        return new CourseInfoDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.HeaderID = source["HeaderID"];
	        this.Name = source["Name"];
	        this.Md5 = source["Md5"];
	        this.Md5s = source["Md5s"];
	        this.Sha256 = source["Sha256"];
	        this.Sha256s = source["Sha256s"];
	        this.NoSepJoinedSha256s = source["NoSepJoinedSha256s"];
	        this.Constranints = source["Constranints"];
	    }
	}
	export class DiffTableDataDto {
	    ID: number;
	    HeaderID: number;
	    Artist: string;
	    Comment: string;
	    Level: string;
	    Lr2BmsId: string;
	    Md5: string;
	    NameDiff: string;
	    Title: string;
	    Url: string;
	    UrlDiff: string;
	    Sha256: string;
	    Lamp: number;
	    GhostLamp: number;
	    PlayCount: number;
	
	    static createFrom(source: any = {}) {
	        return new DiffTableDataDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.HeaderID = source["HeaderID"];
	        this.Artist = source["Artist"];
	        this.Comment = source["Comment"];
	        this.Level = source["Level"];
	        this.Lr2BmsId = source["Lr2BmsId"];
	        this.Md5 = source["Md5"];
	        this.NameDiff = source["NameDiff"];
	        this.Title = source["Title"];
	        this.Url = source["Url"];
	        this.UrlDiff = source["UrlDiff"];
	        this.Sha256 = source["Sha256"];
	        this.Lamp = source["Lamp"];
	        this.GhostLamp = source["GhostLamp"];
	        this.PlayCount = source["PlayCount"];
	    }
	}
	export class DiffTableHeaderDto {
	    ID: number;
	    HeaderUrl: string;
	    DataUrl: string;
	    Name: string;
	    OriginalUrl?: string;
	    Symbol: string;
	    Contents: DiffTableDataDto[];
	    SortedLevels: string[];
	    LevelLayeredContents: Record<string, DiffTableDataDto[]>;
	    Level: string;
	    Children: DiffTableHeaderDto[];
	
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
	        this.Contents = this.convertValues(source["Contents"], DiffTableDataDto);
	        this.SortedLevels = source["SortedLevels"];
	        this.LevelLayeredContents = this.convertValues(source["LevelLayeredContents"], DiffTableDataDto[], true);
	        this.Level = source["Level"];
	        this.Children = this.convertValues(source["Children"], DiffTableHeaderDto);
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
	export class FolderContentDto {
	    ID: number;
	    FolderID: number;
	    FolderName: string;
	    Sha256: string;
	    Md5: string;
	    Title: string;
	
	    static createFrom(source: any = {}) {
	        return new FolderContentDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FolderID = source["FolderID"];
	        this.FolderName = source["FolderName"];
	        this.Sha256 = source["Sha256"];
	        this.Md5 = source["Md5"];
	        this.Title = source["Title"];
	    }
	}
	export class FolderDefinitionDto {
	    name: string;
	    sql: string;
	
	    static createFrom(source: any = {}) {
	        return new FolderDefinitionDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.sql = source["sql"];
	    }
	}
	export class FolderDto {
	    ID: number;
	    FolderName: string;
	    Contents: FolderContentDto[];
	
	    static createFrom(source: any = {}) {
	        return new FolderDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FolderName = source["FolderName"];
	        this.Contents = this.convertValues(source["Contents"], FolderContentDto);
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
	export class RivalScoreLogDto {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    RivalId: number;
	    Sha256: string;
	    Mode: string;
	    Clear: number;
	    OldClear: number;
	    Score: number;
	    OldScore: number;
	    Combo: number;
	    OldCombo: number;
	    Minbp: number;
	    OldMinbp: number;
	    Timestamp: number;
	    RivalSongDataID: number;
	    Title: string;
	    RecordTime: string;
	    Page: number;
	    PageSize: number;
	    PageCount: number;
	
	    static createFrom(source: any = {}) {
	        return new RivalScoreLogDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.RivalId = source["RivalId"];
	        this.Sha256 = source["Sha256"];
	        this.Mode = source["Mode"];
	        this.Clear = source["Clear"];
	        this.OldClear = source["OldClear"];
	        this.Score = source["Score"];
	        this.OldScore = source["OldScore"];
	        this.Combo = source["Combo"];
	        this.OldCombo = source["OldCombo"];
	        this.Minbp = source["Minbp"];
	        this.OldMinbp = source["OldMinbp"];
	        this.Timestamp = source["Timestamp"];
	        this.RivalSongDataID = source["RivalSongDataID"];
	        this.Title = source["Title"];
	        this.RecordTime = source["RecordTime"];
	        this.Page = source["Page"];
	        this.PageSize = source["PageSize"];
	        this.PageCount = source["PageCount"];
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
	export class RivalSongDataDto {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    RivalId: number;
	    Md5: string;
	    Sha256: string;
	    Title: string;
	    SubTitle: string;
	    Genre: string;
	    Artist: string;
	    SubArtist: string;
	    Tag: string;
	    BackBmp: string;
	    Level: number;
	    Difficulty: number;
	    MaxBpm: number;
	    MinBpm: number;
	    Length: number;
	    Mode: number;
	    Judge: number;
	    Feature: number;
	    Content: number;
	    Date: number;
	    Favorite: number;
	    AddDate: number;
	    Notes: number;
	    Page: number;
	    PageSize: number;
	
	    static createFrom(source: any = {}) {
	        return new RivalSongDataDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.RivalId = source["RivalId"];
	        this.Md5 = source["Md5"];
	        this.Sha256 = source["Sha256"];
	        this.Title = source["Title"];
	        this.SubTitle = source["SubTitle"];
	        this.Genre = source["Genre"];
	        this.Artist = source["Artist"];
	        this.SubArtist = source["SubArtist"];
	        this.Tag = source["Tag"];
	        this.BackBmp = source["BackBmp"];
	        this.Level = source["Level"];
	        this.Difficulty = source["Difficulty"];
	        this.MaxBpm = source["MaxBpm"];
	        this.MinBpm = source["MinBpm"];
	        this.Length = source["Length"];
	        this.Mode = source["Mode"];
	        this.Judge = source["Judge"];
	        this.Feature = source["Feature"];
	        this.Content = source["Content"];
	        this.Date = source["Date"];
	        this.Favorite = source["Favorite"];
	        this.AddDate = source["AddDate"];
	        this.Notes = source["Notes"];
	        this.Page = source["Page"];
	        this.PageSize = source["PageSize"];
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
	
	export class CourseInfo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    HeaderID: number;
	    Name: string;
	    Md5s: string;
	    Constranints: string;
	
	    static createFrom(source: any = {}) {
	        return new CourseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.HeaderID = source["HeaderID"];
	        this.Name = source["Name"];
	        this.Md5s = source["Md5s"];
	        this.Constranints = source["Constranints"];
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
	    DataUrl: string;
	    Name: string;
	    OriginalUrl?: string;
	    Symbol: string;
	
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
	        this.DataUrl = source["DataUrl"];
	        this.Name = source["Name"];
	        this.OriginalUrl = source["OriginalUrl"];
	        this.Symbol = source["Symbol"];
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
	export class Folder {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    FolderName: string;
	    BitIndex: number;
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.FolderName = source["FolderName"];
	        this.BitIndex = source["BitIndex"];
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

export namespace vo {
	
	export class RivalScoreLogVo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    RivalId: number;
	    Sha256: string;
	    Mode: string;
	    Clear: number;
	    OldClear: number;
	    Score: number;
	    OldScore: number;
	    Combo: number;
	    OldCombo: number;
	    Minbp: number;
	    OldMinbp: number;
	    Timestamp: number;
	    Page: number;
	    PageSize: number;
	
	    static createFrom(source: any = {}) {
	        return new RivalScoreLogVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.RivalId = source["RivalId"];
	        this.Sha256 = source["Sha256"];
	        this.Mode = source["Mode"];
	        this.Clear = source["Clear"];
	        this.OldClear = source["OldClear"];
	        this.Score = source["Score"];
	        this.OldScore = source["OldScore"];
	        this.Combo = source["Combo"];
	        this.OldCombo = source["OldCombo"];
	        this.Minbp = source["Minbp"];
	        this.OldMinbp = source["OldMinbp"];
	        this.Timestamp = source["Timestamp"];
	        this.Page = source["Page"];
	        this.PageSize = source["PageSize"];
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

