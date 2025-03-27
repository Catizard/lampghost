// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {result} from '../models';
import {vo} from '../models';
import {entity} from '../models';
import {dto} from '../models';

export function AddFolder(arg1:string):Promise<result.RtnMessage>;

export function DelFolder(arg1:number):Promise<result.RtnMessage>;

export function DelFolderContent(arg1:number):Promise<result.RtnMessage>;

export function FindFolderContentList(arg1:vo.FolderContentVo):Promise<result.RtnDataList>;

export function FindFolderList(arg1:vo.FolderVo):Promise<result.RtnDataList>;

export function FindFolderTree(arg1:vo.FolderVo):Promise<result.RtnDataList>;

export function GENERATOR_FOLDER():Promise<entity.Folder>;

export function GENERATOR_FOLDER_CONTENT_DTO():Promise<dto.FolderContentDto>;

export function GENERATOR_FOLDER_DEFINITION_DTO():Promise<dto.FolderDefinitionDto>;

export function GENERATOR_FOLDER_DTO():Promise<dto.FolderDto>;

export function QueryFolderDefinition():Promise<result.RtnDataList>;
