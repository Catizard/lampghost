export namespace config {
	
	export class ApplicationConfig {
	    InternalServerPort: number;
	    FolderSymbol: string;
	    IgnoreVariantCourse: number;
	    Locale: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplicationConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InternalServerPort = source["InternalServerPort"];
	        this.FolderSymbol = source["FolderSymbol"];
	        this.IgnoreVariantCourse = source["IgnoreVariantCourse"];
	        this.Locale = source["Locale"];
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
	    Constraints: string;
	    Clear: number;
	    // Go type: time
	    FirstClearTimestamp: any;
	    Constraint: string[];
	
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
	        this.Constraints = source["Constraints"];
	        this.Clear = source["Clear"];
	        this.FirstClearTimestamp = this.convertValues(source["FirstClearTimestamp"], null);
	        this.Constraint = source["Constraint"];
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
	    PlayCount: number;
	    GhostLamp: number;
	    GhostPlayCount: number;
	    DataLost: boolean;
	
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
	        this.PlayCount = source["PlayCount"];
	        this.GhostLamp = source["GhostLamp"];
	        this.GhostPlayCount = source["GhostPlayCount"];
	        this.DataLost = source["DataLost"];
	    }
	}
	export class DiffTableHeaderDto {
	    ID: number;
	    HeaderUrl: string;
	    DataUrl: string;
	    Name: string;
	    OriginalUrl?: string;
	    Symbol: string;
	    LevelOrders: string;
	    UnjoinedLevelOrder: string[];
	    TagColor: string;
	    TagTextColor: string;
	    NoTagBuild?: number;
	    Contents: DiffTableDataDto[];
	    SortedLevels: string[];
	    LevelLayeredContents: Record<string, DiffTableDataDto[]>;
	    Level: string;
	    Children: DiffTableHeaderDto[];
	    LampCount: Record<number, number>;
	    SongCount: number;
	
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
	        this.LevelOrders = source["LevelOrders"];
	        this.UnjoinedLevelOrder = source["UnjoinedLevelOrder"];
	        this.TagColor = source["TagColor"];
	        this.TagTextColor = source["TagTextColor"];
	        this.NoTagBuild = source["NoTagBuild"];
	        this.Contents = this.convertValues(source["Contents"], DiffTableDataDto);
	        this.SortedLevels = source["SortedLevels"];
	        this.LevelLayeredContents = this.convertValues(source["LevelLayeredContents"], DiffTableDataDto[], true);
	        this.Level = source["Level"];
	        this.Children = this.convertValues(source["Children"], DiffTableHeaderDto);
	        this.LampCount = source["LampCount"];
	        this.SongCount = source["SongCount"];
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
	export class DiffTableTagDto {
	    Md5: string;
	    TableName: string;
	    TableLevel: string;
	    TableSymbol: string;
	    TableTagColor: string;
	    TableTagTextColor: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffTableTagDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Md5 = source["Md5"];
	        this.TableName = source["TableName"];
	        this.TableLevel = source["TableLevel"];
	        this.TableSymbol = source["TableSymbol"];
	        this.TableTagColor = source["TableTagColor"];
	        this.TableTagTextColor = source["TableTagTextColor"];
	    }
	}
	export class FolderContentDto {
	    ID: number;
	    FolderID: number;
	    FolderName: string;
	    Sha256: string;
	    Md5: string;
	    Title: string;
	    Lamp: number;
	
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
	        this.Lamp = source["Lamp"];
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
	export class KeyCountDto {
	    RecordDate: string;
	    KeyCount: number;
	
	    static createFrom(source: any = {}) {
	        return new KeyCountDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RecordDate = source["RecordDate"];
	        this.KeyCount = source["KeyCount"];
	    }
	}
	export class RivalInfoDto {
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
	    ScoreDataLogPath?: string;
	    PlayCount: number;
	    MainUser: boolean;
	    DiffTableHeader?: DiffTableHeaderDto;
	
	    static createFrom(source: any = {}) {
	        return new RivalInfoDto(source);
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
	        this.ScoreDataLogPath = source["ScoreDataLogPath"];
	        this.PlayCount = source["PlayCount"];
	        this.MainUser = source["MainUser"];
	        this.DiffTableHeader = this.convertValues(source["DiffTableHeader"], DiffTableHeaderDto);
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
	    // Go type: time
	    RecordTime: any;
	    Md5: string;
	    RivalSongDataID: number;
	    Title: string;
	    TableTags: DiffTableTagDto[];
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
	        this.RecordTime = this.convertValues(source["RecordTime"], null);
	        this.Md5 = source["Md5"];
	        this.RivalSongDataID = source["RivalSongDataID"];
	        this.Title = source["Title"];
	        this.TableTags = this.convertValues(source["TableTags"], DiffTableTagDto);
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
	export class RivalTagDto {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    RivalId: number;
	    TagName: string;
	    Generated: boolean;
	    Enabled: boolean;
	    // Go type: time
	    RecordTime: any;
	
	    static createFrom(source: any = {}) {
	        return new RivalTagDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.RivalId = source["RivalId"];
	        this.TagName = source["TagName"];
	        this.Generated = source["Generated"];
	        this.Enabled = source["Enabled"];
	        this.RecordTime = this.convertValues(source["RecordTime"], null);
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
	    Sha256s: string;
	    Md5s: string;
	    Constraints: string;
	
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
	        this.Sha256s = source["Sha256s"];
	        this.Md5s = source["Md5s"];
	        this.Constraints = source["Constraints"];
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
	    OrderNumber: number;
	    LevelOrders: string;
	    TagColor: string;
	    TagTextColor: string;
	    NoTagBuild?: number;
	
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
	        this.OrderNumber = source["OrderNumber"];
	        this.LevelOrders = source["LevelOrders"];
	        this.TagColor = source["TagColor"];
	        this.TagTextColor = source["TagTextColor"];
	        this.NoTagBuild = source["NoTagBuild"];
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
	export class Page {
	    page: number;
	    pageSize: number;
	    pageCount: number;
	
	    static createFrom(source: any = {}) {
	        return new Page(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.pageSize = source["pageSize"];
	        this.pageCount = source["pageCount"];
	    }
	}
	export class PredefineTableHeader {
	    HeaderUrl: string;
	    Name: string;
	    Symbol: string;
	    TagColor: string;
	    TagTextColor: string;
	    Category: string;
	
	    static createFrom(source: any = {}) {
	        return new PredefineTableHeader(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HeaderUrl = source["HeaderUrl"];
	        this.Name = source["Name"];
	        this.Symbol = source["Symbol"];
	        this.TagColor = source["TagColor"];
	        this.TagTextColor = source["TagTextColor"];
	        this.Category = source["Category"];
	    }
	}
	export class PredefineTableScheme {
	    Headers: PredefineTableHeader[];
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new PredefineTableScheme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Headers = this.convertValues(source["Headers"], PredefineTableHeader);
	        this.Name = source["Name"];
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
	    ScoreDataLogPath?: string;
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
	        this.ScoreDataLogPath = source["ScoreDataLogPath"];
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
	export class RivalTag {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    RivalId: number;
	    TagName: string;
	    Generated: boolean;
	    Enabled: boolean;
	    // Go type: time
	    RecordTime: any;
	
	    static createFrom(source: any = {}) {
	        return new RivalTag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.RivalId = source["RivalId"];
	        this.TagName = source["TagName"];
	        this.Generated = source["Generated"];
	        this.Enabled = source["Enabled"];
	        this.RecordTime = this.convertValues(source["RecordTime"], null);
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
	export class RtnPage {
	    Pagination: entity.Page;
	    Rows: any[];
	    Code: number;
	    Msg: string;
	    // Go type: time
	    Timestamp: any;
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new RtnPage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Pagination = this.convertValues(source["Pagination"], entity.Page);
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

}

export namespace vo {
	
	export class ChartInfoVo {
	    Title: string;
	    SubTitle: string;
	    Artist: string;
	    sha256: string;
	    md5: string;
	
	    static createFrom(source: any = {}) {
	        return new ChartInfoVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.SubTitle = source["SubTitle"];
	        this.Artist = source["Artist"];
	        this.sha256 = source["sha256"];
	        this.md5 = source["md5"];
	    }
	}
	export class CourseInfoVo {
	    name: string;
	    md5: string[];
	    sha256: string[];
	    constraint: string[];
	    charts: ChartInfoVo[];
	    HeaderID: number;
	
	    static createFrom(source: any = {}) {
	        return new CourseInfoVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.md5 = source["md5"];
	        this.sha256 = source["sha256"];
	        this.constraint = source["constraint"];
	        this.charts = this.convertValues(source["charts"], ChartInfoVo);
	        this.HeaderID = source["HeaderID"];
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
	export class DiffTableHeaderVo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    data_url: string;
	    name: string;
	    original_url?: string;
	    symbol: string;
	    Courses: CourseInfoVo[];
	    HeaderUrl: string;
	    LevelOrders: string;
	    level_order: string[];
	    TagColor: string;
	    TagTextColor: string;
	    NoTagBuild?: number;
	    Level: string;
	    RivalID: number;
	    GhostRivalID: number;
	    GhostRivalTagID: number;
	    Pagination?: entity.Page;
	    SortBy?: string;
	    SortOrder?: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffTableHeaderVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.data_url = source["data_url"];
	        this.name = source["name"];
	        this.original_url = source["original_url"];
	        this.symbol = source["symbol"];
	        this.Courses = this.convertValues(source["Courses"], CourseInfoVo);
	        this.HeaderUrl = source["HeaderUrl"];
	        this.LevelOrders = source["LevelOrders"];
	        this.level_order = source["level_order"];
	        this.TagColor = source["TagColor"];
	        this.TagTextColor = source["TagTextColor"];
	        this.NoTagBuild = source["NoTagBuild"];
	        this.Level = source["Level"];
	        this.RivalID = source["RivalID"];
	        this.GhostRivalID = source["GhostRivalID"];
	        this.GhostRivalTagID = source["GhostRivalTagID"];
	        this.Pagination = this.convertValues(source["Pagination"], entity.Page);
	        this.SortBy = source["SortBy"];
	        this.SortOrder = source["SortOrder"];
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
	export class FolderContentVo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    FolderID: number;
	    FolderName: string;
	    Sha256: string;
	    Md5: string;
	    Title: string;
	    IDs: number[];
	    FolderIDs: number[];
	
	    static createFrom(source: any = {}) {
	        return new FolderContentVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.FolderID = source["FolderID"];
	        this.FolderName = source["FolderName"];
	        this.Sha256 = source["Sha256"];
	        this.Md5 = source["Md5"];
	        this.Title = source["Title"];
	        this.IDs = source["IDs"];
	        this.FolderIDs = source["FolderIDs"];
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
	export class FolderVo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    FolderName: string;
	    BitIndex: number;
	    IDs: number[];
	    IgnoreSha256?: string;
	    IgnoreRivalSongDataID?: number;
	    RivalID: number;
	
	    static createFrom(source: any = {}) {
	        return new FolderVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.FolderName = source["FolderName"];
	        this.BitIndex = source["BitIndex"];
	        this.IDs = source["IDs"];
	        this.IgnoreSha256 = source["IgnoreSha256"];
	        this.IgnoreRivalSongDataID = source["IgnoreRivalSongDataID"];
	        this.RivalID = source["RivalID"];
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
	export class RivalInfoVo {
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
	    ScoreDataLogPath?: string;
	    PlayCount: number;
	    MainUser: boolean;
	    Pagination?: entity.Page;
	    Locale?: string;
	
	    static createFrom(source: any = {}) {
	        return new RivalInfoVo(source);
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
	        this.ScoreDataLogPath = source["ScoreDataLogPath"];
	        this.PlayCount = source["PlayCount"];
	        this.MainUser = source["MainUser"];
	        this.Pagination = this.convertValues(source["Pagination"], entity.Page);
	        this.Locale = source["Locale"];
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
	export class RivalScoreDataLogVo {
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
	    // Go type: time
	    RecordTime: any;
	    Epg: number;
	    Lpg: number;
	    Egr: number;
	    Lgr: number;
	    Egd: number;
	    Lgd: number;
	    Ebd: number;
	    Lbd: number;
	    Epr: number;
	    Lpr: number;
	    Ems: number;
	    Lms: number;
	    Notes: number;
	    Combo: number;
	    Minbp: number;
	    PlayCount: number;
	    ClearCount: number;
	    Option: number;
	    Seed: number;
	    Random: number;
	    State: number;
	    SpecifyYear?: string;
	
	    static createFrom(source: any = {}) {
	        return new RivalScoreDataLogVo(source);
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
	        this.RecordTime = this.convertValues(source["RecordTime"], null);
	        this.Epg = source["Epg"];
	        this.Lpg = source["Lpg"];
	        this.Egr = source["Egr"];
	        this.Lgr = source["Lgr"];
	        this.Egd = source["Egd"];
	        this.Lgd = source["Lgd"];
	        this.Ebd = source["Ebd"];
	        this.Lbd = source["Lbd"];
	        this.Epr = source["Epr"];
	        this.Lpr = source["Lpr"];
	        this.Ems = source["Ems"];
	        this.Lms = source["Lms"];
	        this.Notes = source["Notes"];
	        this.Combo = source["Combo"];
	        this.Minbp = source["Minbp"];
	        this.PlayCount = source["PlayCount"];
	        this.ClearCount = source["ClearCount"];
	        this.Option = source["Option"];
	        this.Seed = source["Seed"];
	        this.Random = source["Random"];
	        this.State = source["State"];
	        this.SpecifyYear = source["SpecifyYear"];
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
	    // Go type: time
	    RecordTime: any;
	    Pagination?: entity.Page;
	    OnlyCourseLogs: boolean;
	    NoCourseLog: boolean;
	    // Go type: time
	    StartRecordTime: any;
	    // Go type: time
	    EndRecordTime: any;
	    StartRecordTimestamp: number;
	    EndRecordTimestamp: number;
	    MinimumClear?: number;
	    SpecifyYear?: string;
	    SongNameLike?: string;
	    Sha256s: string[];
	    HeaderID: number;
	
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
	        this.RecordTime = this.convertValues(source["RecordTime"], null);
	        this.Pagination = this.convertValues(source["Pagination"], entity.Page);
	        this.OnlyCourseLogs = source["OnlyCourseLogs"];
	        this.NoCourseLog = source["NoCourseLog"];
	        this.StartRecordTime = this.convertValues(source["StartRecordTime"], null);
	        this.EndRecordTime = this.convertValues(source["EndRecordTime"], null);
	        this.StartRecordTimestamp = source["StartRecordTimestamp"];
	        this.EndRecordTimestamp = source["EndRecordTimestamp"];
	        this.MinimumClear = source["MinimumClear"];
	        this.SpecifyYear = source["SpecifyYear"];
	        this.SongNameLike = source["SongNameLike"];
	        this.Sha256s = source["Sha256s"];
	        this.HeaderID = source["HeaderID"];
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
	export class RivalTagVo {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    RivalId: number;
	    TagName: string;
	    Generated: boolean;
	    Enabled: boolean;
	    // Go type: time
	    RecordTime: any;
	    Pagination?: entity.Page;
	    RecordTimestamp?: number;
	
	    static createFrom(source: any = {}) {
	        return new RivalTagVo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.RivalId = source["RivalId"];
	        this.TagName = source["TagName"];
	        this.Generated = source["Generated"];
	        this.Enabled = source["Enabled"];
	        this.RecordTime = this.convertValues(source["RecordTime"], null);
	        this.Pagination = this.convertValues(source["Pagination"], entity.Page);
	        this.RecordTimestamp = source["RecordTimestamp"];
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

