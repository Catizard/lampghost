// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {vo} from '../models';
import {result} from '../models';
import {entity} from '../models';
import {dto} from '../models';
import {download} from '../models';
import {config} from '../models';

export function AddBatchDiffTableHeader(arg1:Array<vo.DiffTableHeaderVo>):Promise<result.RtnDataList>;

export function AddCustomCourse(arg1:vo.CustomCourseVo):Promise<result.RtnMessage>;

export function AddCustomCourseData(arg1:entity.CustomCourseData):Promise<result.RtnMessage>;

export function AddCustomDiffTable(arg1:vo.CustomDiffTableVo):Promise<result.RtnMessage>;

export function AddDiffTableHeader(arg1:vo.DiffTableHeaderVo):Promise<result.RtnMessage>;

export function AddFolder(arg1:vo.FolderVo):Promise<result.RtnMessage>;

export function AddRivalInfo(arg1:vo.RivalInfoVo):Promise<result.RtnMessage>;

export function AddRivalTag(arg1:vo.RivalTagVo):Promise<result.RtnMessage>;

export function BindFolderContentToCustomCourse(arg1:number,arg2:number):Promise<result.RtnMessage>;

export function BindSongToCustomCourse(arg1:string,arg2:string,arg3:number):Promise<result.RtnMessage>;

export function BindSongToFolder(arg1:vo.BindToFolderVo):Promise<result.RtnMessage>;

export function CancelDownloadTask(arg1:number):Promise<result.RtnMessage>;

export function ChooseBeatorajaDirectory():Promise<result.RtnData>;

export function DelDiffTableHeader(arg1:number):Promise<result.RtnMessage>;

export function DelFolder(arg1:number):Promise<result.RtnMessage>;

export function DelFolderContent(arg1:number):Promise<result.RtnMessage>;

export function DelRivalInfo(arg1:number):Promise<result.RtnMessage>;

export function DeleteCustomCourse(arg1:number):Promise<result.RtnMessage>;

export function DeleteCustomCourseData(arg1:number):Promise<result.RtnMessage>;

export function DeleteCustomDiffTable(arg1:number):Promise<result.RtnMessage>;

export function DeleteRivalTagByID(arg1:number):Promise<result.RtnMessage>;

export function FindCourseInfoList():Promise<result.RtnDataList>;

export function FindCourseInfoListWithRival(arg1:vo.CourseInfoVo):Promise<result.RtnDataList>;

export function FindCustomCourseByID(arg1:number):Promise<result.RtnData>;

export function FindCustomCourseList(arg1:vo.CustomCourseVo):Promise<result.RtnDataList>;

export function FindCustomDiffTableByID(arg1:number):Promise<result.RtnData>;

export function FindCustomDiffTableList(arg1:vo.CustomDiffTableVo):Promise<result.RtnDataList>;

export function FindDiffTableHeaderList():Promise<result.RtnDataList>;

export function FindDiffTableHeaderTree(arg1:vo.DiffTableHeaderVo):Promise<result.RtnDataList>;

export function FindDiffTableHeaderTreeWithRival(arg1:vo.DiffTableHeaderVo):Promise<result.RtnDataList>;

export function FindDownloadTaskList():Promise<result.RtnDataList>;

export function FindFolderContentList(arg1:vo.FolderContentVo):Promise<result.RtnDataList>;

export function FindFolderList(arg1:vo.FolderVo):Promise<result.RtnDataList>;

export function FindFolderTree(arg1:vo.FolderVo):Promise<result.RtnDataList>;

export function FindRivalInfoList():Promise<result.RtnDataList>;

export function FindRivalTagByID(arg1:number):Promise<result.RtnData>;

export function FindRivalTagList(arg1:vo.RivalTagVo):Promise<result.RtnDataList>;

export function GENERATE_RIVAL_SCORE_LOG():Promise<dto.RivalScoreLogDto>;

export function GENERATE_RIVAL_SONG_DATA_DTO():Promise<dto.RivalSongDataDto>;

export function GENERATE_RIVAL_TAG():Promise<entity.RivalTag>;

export function GENERATE_RIVAL_TAG_DTO():Promise<dto.RivalTagDto>;

export function GENERATOR_BEATORAJA_DIRECTORY_META():Promise<dto.BeatorajaDirectoryMeta>;

export function GENERATOR_COURSE_INFO():Promise<entity.CourseInfo>;

export function GENERATOR_COURSE_INFO_DTO():Promise<dto.CourseInfoDto>;

export function GENERATOR_CUSTOM_COUSE():Promise<entity.CustomCourse>;

export function GENERATOR_CUSTOM_COUSE_DATA():Promise<dto.CustomCourseDataDto>;

export function GENERATOR_CUSTOM_COUSE_DTO():Promise<dto.CustomCourseDto>;

export function GENERATOR_CUSTOM_DIFF_TABLE_DTO():Promise<dto.CustomDiffTableDto>;

export function GENERATOR_DOWNLOAD_SOURCE():Promise<download.DownloadSource>;

export function GENERATOR_DOWNLOAD_TASK():Promise<entity.DownloadTask>;

export function GENERATOR_FOLDER():Promise<entity.Folder>;

export function GENERATOR_FOLDER_CONTENT_DTO():Promise<dto.FolderContentDto>;

export function GENERATOR_FOLDER_DTO():Promise<dto.FolderDto>;

export function GENERATOR_KEY_COUNT_DTO():Promise<dto.KeyCountDto>;

export function GENERATOR_NOTIFICATION_DTO():Promise<dto.NotificationDto>;

export function GENERATOR_PREDEFINE_TABLE_HEADER():Promise<entity.PredefineTableHeader>;

export function GENERATOR_PREDEFINE_TABLE_SCHEME():Promise<entity.PredefineTableScheme>;

export function GENERATOR_RIVAL_INFO():Promise<entity.RivalInfo>;

export function GENERATOR_RIVAL_INFO_DTO():Promise<dto.RivalInfoDto>;

export function GENERATOR_TABLE_DATA():Promise<dto.DiffTableDataDto>;

export function GENERATOR_TABLE_HEADER():Promise<entity.DiffTableHeader>;

export function GENERATOR_TABLE_HEADER_DTO():Promise<dto.DiffTableHeaderDto>;

export function GENERATOR_TABLE_TAG_DTO():Promise<dto.DiffTableTagDto>;

export function InitializeMainUser(arg1:vo.InitializeRivalInfoVo):Promise<result.RtnMessage>;

export function OpenDirectoryDialog(arg1:string):Promise<result.RtnData>;

export function OpenFileDialog(arg1:string):Promise<result.RtnData>;

export function QueryCourseSongListWithRival(arg1:vo.CourseInfoVo):Promise<result.RtnDataList>;

export function QueryCurrentDownloadSource():Promise<result.RtnData>;

export function QueryCustomCourseSongListWithRival(arg1:vo.CustomCourseVo):Promise<result.RtnDataList>;

export function QueryCustomDiffTablePageList(arg1:vo.CustomDiffTableVo):Promise<result.RtnPage>;

export function QueryDiffTableDataWithRival(arg1:vo.DiffTableHeaderVo):Promise<result.RtnPage>;

export function QueryDiffTableInfoById(arg1:number):Promise<result.RtnData>;

export function QueryFolderContentWithRival(arg1:vo.FolderContentVo):Promise<result.RtnPage>;

export function QueryLatestVersion():Promise<result.RtnMessage>;

export function QueryMainUser():Promise<result.RtnData>;

export function QueryPredefineTableSchemes():Promise<result.RtnDataList>;

export function QueryPrevDayScoreLogList(arg1:vo.RivalScoreLogVo):Promise<result.RtnDataList>;

export function QueryPreviewURLByMd5(arg1:string):Promise<result.RtnData>;

export function QueryRivalInfoPageList(arg1:vo.RivalInfoVo):Promise<result.RtnPage>;

export function QueryRivalPlayedYears(arg1:number):Promise<result.RtnDataList>;

export function QueryRivalScoreLogPageList(arg1:vo.RivalScoreLogVo):Promise<result.RtnPage>;

export function QueryRivalTagPageList(arg1:vo.RivalTagVo):Promise<result.RtnPage>;

export function QuerySongDataPageList(arg1:vo.RivalSongDataVo):Promise<result.RtnPage>;

export function QueryUserInfoByID(arg1:number):Promise<result.RtnData>;

export function QueryUserInfoWithLevelLayeredDiffTableLampStatus(arg1:number,arg2:number):Promise<result.RtnData>;

export function QueryUserKeyCountInYear(arg1:vo.RivalScoreDataLogVo):Promise<result.RtnDataList>;

export function QueryUserPlayCountInYear(arg1:number,arg2:string):Promise<result.RtnDataList>;

export function ReadConfig():Promise<result.RtnData>;

export function ReloadDiffTableHeader(arg1:number):Promise<result.RtnMessage>;

export function ReloadRivalData(arg1:number,arg2:boolean):Promise<result.RtnMessage>;

export function ReloadRivalSongData():Promise<result.RtnMessage>;

export function RestartDownloadTask(arg1:number):Promise<result.RtnMessage>;

export function RunServer():Promise<void>;

export function SetScoreLogFilePath(arg1:string):Promise<void>;

export function SubmitSingleMD5DownloadTask(arg1:string,arg2:any):Promise<result.RtnMessage>;

export function SupplyMissingBMSFromTable(arg1:number):Promise<result.RtnMessage>;

export function SyncRivalTag(arg1:number):Promise<result.RtnMessage>;

export function UpdateCustomCourse(arg1:vo.CustomCourseVo):Promise<result.RtnMessage>;

export function UpdateCustomCourseDataOrder(arg1:Array<number>):Promise<result.RtnMessage>;

export function UpdateCustomCourseOrder(arg1:Array<number>):Promise<result.RtnMessage>;

export function UpdateCustomDiffTable(arg1:vo.CustomDiffTableVo):Promise<result.RtnMessage>;

export function UpdateDiffTableHeader(arg1:vo.DiffTableHeaderVo):Promise<result.RtnMessage>;

export function UpdateFolderOrder(arg1:Array<number>):Promise<result.RtnMessage>;

export function UpdateHeaderLevelOrders(arg1:vo.DiffTableHeaderVo):Promise<result.RtnMessage>;

export function UpdateHeaderOrder(arg1:Array<number>):Promise<result.RtnMessage>;

export function UpdateRivalInfo(arg1:vo.RivalInfoVo):Promise<result.RtnMessage>;

export function UpdateRivalReverseImportInfo(arg1:vo.RivalInfoVo):Promise<result.RtnMessage>;

export function UpdateRivalTag(arg1:vo.RivalTagUpdateParam):Promise<result.RtnMessage>;

export function WriteConfig(arg1:config.ApplicationConfig):Promise<result.RtnMessage>;
