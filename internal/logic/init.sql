
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `fms_files` (
  `id` char(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID',
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `status` tinyint unsigned DEFAULT '1' COMMENT 'Status 1: normal 2: ban | 状态 1 正常 2 禁用',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'File''s name | 文件名称',
  `file_type` tinyint unsigned NOT NULL COMMENT 'File''s type | 文件类型',
  `size` bigint unsigned NOT NULL COMMENT 'File''s size | 文件大小',
  `path` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'File''s path | 文件路径',
  `user_id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'User''s UUID | 用户的 UUID',
  `md5` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The md5 of the file | 文件的 md5',
  PRIMARY KEY (`id`),
  KEY `file_user_id` (`user_id`),
  KEY `file_file_type` (`file_type`),
  KEY `file_path` (`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='File Table | 文件表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `fms_file_tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `status` tinyint unsigned DEFAULT '1' COMMENT 'Status 1: normal 2: ban | 状态 1 正常 2 禁用',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'FileTag''s name | 标签名称',
  `remark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'The remark of tag | 标签的备注',
  PRIMARY KEY (`id`),
  KEY `filetag_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='File''s Tags Table | 文件标签表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `fms_storage_providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `state` tinyint(1) DEFAULT '1' COMMENT 'State true: normal false: ban | 状态 true 正常 false 禁用',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The cloud storage service name | 服务名称',
  `bucket` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The cloud storage bucket name | 云存储服务的存储桶',
  `secret_id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The secret ID | 密钥 ID',
  `secret_key` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The secret key | 密钥 Key',
  `endpoint` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The service URL | 服务器地址',
  `folder` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'The folder in cloud | 云服务目标文件夹',
  `region` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The service region | 服务器所在地区',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Is it the default provider | 是否为默认提供商',
  `use_cdn` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Does it use CDN | 是否使用 CDN',
  `cdn_url` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'CDN URL | CDN 地址',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Cloud File''s Provider Configuration Table | 云文件服务商配置表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `fms_cloud_files` (
  `id` char(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID',
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `state` tinyint(1) DEFAULT '1' COMMENT 'State true: normal false: ban | 状态 true 正常 false 禁用',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The file''s name | 文件名',
  `url` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The file''s url | 文件地址',
  `size` bigint unsigned NOT NULL COMMENT 'The file''s size | 文件大小',
  `file_type` tinyint unsigned NOT NULL COMMENT 'The file''s type | 文件类型',
  `user_id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The user who upload the file | 上传用户的 ID',
  `cloud_file_storage_providers` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cloudfile_name` (`name`),
  KEY `cloudfile_file_type` (`file_type`),
  KEY `cloudfile_url` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Cloud File Table | 云文件表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `fms_cloud_file_tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `status` tinyint unsigned DEFAULT '1' COMMENT 'Status 1: normal 2: ban | 状态 1 正常 2 禁用',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'CloudFileTag''s name | 标签名称',
  `remark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'The remark of tag | 标签的备注',
  PRIMARY KEY (`id`),
  KEY `cloudfiletag_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Cloud File''s Tags Table | 云文件标签表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `file_tag_files` (
  `file_tag_id` bigint unsigned NOT NULL,
  `file_id` char(36) COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`file_tag_id`,`file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `cloud_file_tag_cloud_files` (
  `cloud_file_tag_id` bigint unsigned NOT NULL,
  `cloud_file_id` char(36) COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`cloud_file_tag_id`,`cloud_file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

-- Seed FMS APIs, casbin policies and menus. INSERT IGNORE keeps this idempotent.
INSERT IGNORE INTO `sys_apis` VALUES (2001,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/init/database','apiDesc.initDatabase','file','Fms','GET',0),(2002,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/upload','apiDesc.uploadFile','file','Fms','POST',0),(2003,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file/list','apiDesc.fileList','file','Fms','POST',0),(2004,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file/update','apiDesc.updateFileInfo','file','Fms','POST',0),(2005,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file/delete','apiDesc.deleteFile','file','Fms','POST',0),(2006,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file/status','apiDesc.setPublicStatus','file','Fms','POST',0),(2007,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file/delete_by_url','apiDesc.deleteFile','file','Fms','POST',0),(2008,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file/download/:id','apiDesc.downloadFile','file','Fms','GET',0),(2009,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file_tag/create','apiDesc.createFileTag','file_tag','Fms','POST',0),(2010,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file_tag/update','apiDesc.updateFileTag','file_tag','Fms','POST',0),(2011,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file_tag/delete','apiDesc.deleteFileTag','file_tag','Fms','POST',0),(2012,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file_tag/list','apiDesc.getFileTagList','file_tag','Fms','POST',0),(2013,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/file_tag','apiDesc.getFileTagById','file_tag','Fms','POST',0),(2014,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/storage_provider/create','apiDesc.createStorageProvider','storage_provider','Fms','POST',0),(2015,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/storage_provider/update','apiDesc.updateStorageProvider','storage_provider','Fms','POST',0),(2016,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/storage_provider/delete','apiDesc.deleteStorageProvider','storage_provider','Fms','POST',0),(2017,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/storage_provider/list','apiDesc.getStorageProviderList','storage_provider','Fms','POST',0),(2018,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/storage_provider','apiDesc.getStorageProviderById','storage_provider','Fms','POST',0),(2019,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file/create','apiDesc.createCloudFile','cloud_file','Fms','POST',0),(2020,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file/update','apiDesc.updateCloudFile','cloud_file','Fms','POST',0),(2021,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file/delete','apiDesc.deleteCloudFile','cloud_file','Fms','POST',0),(2022,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file/list','apiDesc.getCloudFileList','cloud_file','Fms','POST',0),(2023,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file','apiDesc.getCloudFileById','cloud_file','Fms','POST',0),(2024,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file/delete_by_url','apiDesc.deleteFile','cloud_file','Fms','POST',0),(2025,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file/upload','apiDesc.uploadFileToCloud','cloud_file','Fms','POST',0),(2026,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file_tag/create','apiDesc.createCloudFileTag','cloud_file_tag','Fms','POST',0),(2027,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file_tag/update','apiDesc.updateCloudFileTag','cloud_file_tag','Fms','POST',0),(2028,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file_tag/delete','apiDesc.deleteCloudFileTag','cloud_file_tag','Fms','POST',0),(2029,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file_tag/list','apiDesc.getCloudFileTagList','cloud_file_tag','Fms','POST',0),(2030,'2026-10-17 00:00:00','2026-10-17 00:00:00','/fms-api/cloud_file_tag','apiDesc.getCloudFileTagById','cloud_file_tag','Fms','POST',0);

INSERT IGNORE INTO `casbin_rules` VALUES (7001,'p','001','/fms-api/init/database','GET','','',''),(7002,'p','001','/fms-api/upload','POST','','',''),(7003,'p','001','/fms-api/file/list','POST','','',''),(7004,'p','001','/fms-api/file/update','POST','','',''),(7005,'p','001','/fms-api/file/delete','POST','','',''),(7006,'p','001','/fms-api/file/status','POST','','',''),(7007,'p','001','/fms-api/file/delete_by_url','POST','','',''),(7008,'p','001','/fms-api/file/download/:id','GET','','',''),(7009,'p','001','/fms-api/file_tag/create','POST','','',''),(7010,'p','001','/fms-api/file_tag/update','POST','','',''),(7011,'p','001','/fms-api/file_tag/delete','POST','','',''),(7012,'p','001','/fms-api/file_tag/list','POST','','',''),(7013,'p','001','/fms-api/file_tag','POST','','',''),(7014,'p','001','/fms-api/storage_provider/create','POST','','',''),(7015,'p','001','/fms-api/storage_provider/update','POST','','',''),(7016,'p','001','/fms-api/storage_provider/delete','POST','','',''),(7017,'p','001','/fms-api/storage_provider/list','POST','','',''),(7018,'p','001','/fms-api/storage_provider','POST','','',''),(7019,'p','001','/fms-api/cloud_file/create','POST','','',''),(7020,'p','001','/fms-api/cloud_file/update','POST','','',''),(7021,'p','001','/fms-api/cloud_file/delete','POST','','',''),(7022,'p','001','/fms-api/cloud_file/list','POST','','',''),(7023,'p','001','/fms-api/cloud_file','POST','','',''),(7024,'p','001','/fms-api/cloud_file/delete_by_url','POST','','',''),(7025,'p','001','/fms-api/cloud_file/upload','POST','','',''),(7026,'p','001','/fms-api/cloud_file_tag/create','POST','','',''),(7027,'p','001','/fms-api/cloud_file_tag/update','POST','','',''),(7028,'p','001','/fms-api/cloud_file_tag/delete','POST','','',''),(7029,'p','001','/fms-api/cloud_file_tag/list','POST','','',''),(7030,'p','001','/fms-api/cloud_file_tag','POST','','','');

INSERT IGNORE INTO `sys_menus` VALUES (200,'2026-10-17 00:00:00','2026-10-17 00:00:00',3,1,1,'/fms_dir','FileManagementDirectory','','LAYOUT',0,'Fms','','route.fileManagement','ant-design:folder-open-outlined',0,0,0,0,'',0,0,0,20,'',1000000),(201,'2026-10-17 00:00:00','2026-10-17 00:00:00',1,2,1,'/fms/file','FileManagement','','/fms/file/index',0,'Fms',NULL,'route.fileManagement','ant-design:folder-open-outlined',0,0,0,0,'',0,0,0,20,'',200),(202,'2026-10-17 00:00:00','2026-10-17 00:00:00',2,2,1,'/fms/file_tag','FileTagManagement','','/fms/fileTag/index',0,'Fms',NULL,'route.fileTagManagement','ant-design:book-outlined',0,0,0,0,'',0,0,0,20,'',200),(203,'2026-10-17 00:00:00','2026-10-17 00:00:00',3,2,1,'/fms/storage_provider','StorageProviderManagement','','/fms/storageProvider/index',0,'Fms',NULL,'route.storageProviderManagement','mdi:cloud-outline',0,0,0,0,'',0,0,0,20,'',200),(204,'2026-10-17 00:00:00','2026-10-17 00:00:00',4,2,1,'/fms/cloud_file','CloudFileManagement','','/fms/cloudFile/index',0,'Fms',NULL,'route.cloudFileManagement','ant-design:folder-open-outlined',0,0,0,0,'',0,0,0,20,'',200),(205,'2026-10-17 00:00:00','2026-10-17 00:00:00',5,2,1,'/fms/cloud_file_tag','CloudFileTagManagement','','/fms/cloudFileTag/index',0,'Fms',NULL,'route.cloudFileTagManagement','ant-design:book-outlined',0,0,0,0,'',0,0,0,20,'',200);

INSERT IGNORE INTO `role_menus` VALUES (1,200),(1,201),(1,202),(1,203),(1,204),(1,205);
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

