package appserver

const (
	sql_createtable1 = `
	CREATE TABLE [dbo].[iFixsvr_JF_Devices] (
		[ID] [decimal](18, 0) IDENTITY (1, 1) NOT NULL ,
		[DIP] [varchar] (255) COLLATE Chinese_PRC_CI_AS NOT NULL ,
		[DModTime] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[DDesc] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL 
	) ON [PRIMARY]
	`

	sql_createtable2 = `
	CREATE TABLE [dbo].[iFixsvr_JF_Files] (
		[ID] [decimal](18, 0) IDENTITY (1, 1) NOT NULL ,
		[FID] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[FilePath] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[FMd5] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[FModTime] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[FFinish] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[FProgress] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[FDesc] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL 
	) ON [PRIMARY]
	`

	sql_createtable3 = `
	CREATE TABLE [dbo].[iFixsvr_JF_Info] (
		[ID] [decimal](18, 0) IDENTITY (1, 1) NOT NULL ,
		[IID] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField1] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField2] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField3] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField4] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField5] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField6] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField7] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IField8] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IModTime] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL ,
		[IDesc] [varchar] (255) COLLATE Chinese_PRC_CI_AS NULL 
	) ON [PRIMARY]
	`

	sql_devices = `
	IF EXISTS (SELECT * FROM [iFixsvr_JF_Devices] WHERE DIP= '%s' ) 
	UPDATE [iFixsvr_JF_Devices] SET DModTime= '%s' WHERE DIP= '%s' 
	ELSE
	 INSERT INTO [iFixsvr_JF_Devices] (DIP,DModTime) VALUES ('%s','%s')
	 `

	sql_fileexists = `
	 SELECT ID FROM [iFixsvr_JF_Files] WHERE FID = '%s' and FilePath ='%s'
	 `

	sql_insertfile = `
	 INSERT INTO [iFixsvr_JF_Files] (FID,FilePath,FMd5,FModTime,FFinish) VALUES ('%s','%s','%s','%s','%s')
	 `
)
