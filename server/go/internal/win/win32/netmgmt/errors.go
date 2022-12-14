//go:build windows && amd64

package netmgmt

const NERR_Success = 0

const NERR_BASE = 2100

const NERR_NetNotStarted = (NERR_BASE + 2)
const NERR_UnknownServer = (NERR_BASE + 3)
const NERR_ShareMem = (NERR_BASE + 4)
const NERR_NoNetworkResource = (NERR_BASE + 5)
const NERR_RemoteOnly = (NERR_BASE + 6)
const NERR_DevNotRedirected = (NERR_BASE + 7)
const NERR_ServerNotStarted = (NERR_BASE + 14)
const NERR_ItemNotFound = (NERR_BASE + 15)
const NERR_UnknownDevDir = (NERR_BASE + 16)
const NERR_RedirectedPath = (NERR_BASE + 17)
const NERR_DuplicateShare = (NERR_BASE + 18)
const NERR_NoRoom = (NERR_BASE + 19)
const NERR_TooManyItems = (NERR_BASE + 21)
const NERR_InvalidMaxUsers = (NERR_BASE + 22)
const NERR_BufTooSmall = (NERR_BASE + 23)
const NERR_RemoteErr = (NERR_BASE + 27)
const NERR_LanmanIniError = (NERR_BASE + 31)
const NERR_NetworkError = (NERR_BASE + 36)
const NERR_WkstaInconsistentState = (NERR_BASE + 37)
const NERR_WkstaNotStarted = (NERR_BASE + 38)
const NERR_BrowserNotStarted = (NERR_BASE + 39)
const NERR_InternalError = (NERR_BASE + 40)
const NERR_BadTransactConfig = (NERR_BASE + 41)
const NERR_InvalidAPI = (NERR_BASE + 42)
const NERR_BadEventName = (NERR_BASE + 43)
const NERR_DupNameReboot = (NERR_BASE + 44)
const NERR_CfgCompNotFound = (NERR_BASE + 46)
const NERR_CfgParamNotFound = (NERR_BASE + 47)
const NERR_LineTooLong = (NERR_BASE + 49)
const NERR_QNotFound = (NERR_BASE + 50)
const NERR_JobNotFound = (NERR_BASE + 51)
const NERR_DestNotFound = (NERR_BASE + 52)
const NERR_DestExists = (NERR_BASE + 53)
const NERR_QExists = (NERR_BASE + 54)
const NERR_QNoRoom = (NERR_BASE + 55)
const NERR_JobNoRoom = (NERR_BASE + 56)
const NERR_DestNoRoom = (NERR_BASE + 57)
const NERR_DestIdle = (NERR_BASE + 58)
const NERR_DestInvalidOp = (NERR_BASE + 59)
const NERR_ProcNoRespond = (NERR_BASE + 60)
const NERR_SpoolerNotLoaded = (NERR_BASE + 61)
const NERR_DestInvalidState = (NERR_BASE + 62)
const NERR_QInvalidState = (NERR_BASE + 63)
const NERR_JobInvalidState = (NERR_BASE + 64)
const NERR_SpoolNoMemory = (NERR_BASE + 65)
const NERR_DriverNotFound = (NERR_BASE + 66)
const NERR_DataTypeInvalid = (NERR_BASE + 67)
const NERR_ProcNotFound = (NERR_BASE + 68)
const NERR_ServiceTableLocked = (NERR_BASE + 80)
const NERR_ServiceTableFull = (NERR_BASE + 81)
const NERR_ServiceInstalled = (NERR_BASE + 82)
const NERR_ServiceEntryLocked = (NERR_BASE + 83)
const NERR_ServiceNotInstalled = (NERR_BASE + 84)
const NERR_BadServiceName = (NERR_BASE + 85)
const NERR_ServiceCtlTimeout = (NERR_BASE + 86)
const NERR_ServiceCtlBusy = (NERR_BASE + 87)
const NERR_BadServiceProgName = (NERR_BASE + 88)
const NERR_ServiceNotCtrl = (NERR_BASE + 89)
const NERR_ServiceKillProc = (NERR_BASE + 90)
const NERR_ServiceCtlNotValid = (NERR_BASE + 91)
const NERR_NotInDispatchTbl = (NERR_BASE + 92)
const NERR_BadControlRecv = (NERR_BASE + 93)
const NERR_ServiceNotStarting = (NERR_BASE + 94)
const NERR_AlreadyLoggedOn = (NERR_BASE + 100)
const NERR_NotLoggedOn = (NERR_BASE + 101)
const NERR_BadUsername = (NERR_BASE + 102)
const NERR_BadPassword = (NERR_BASE + 103)
const NERR_UnableToAddName_W = (NERR_BASE + 104)
const NERR_UnableToAddName_F = (NERR_BASE + 105)
const NERR_UnableToDelName_W = (NERR_BASE + 106)
const NERR_UnableToDelName_F = (NERR_BASE + 107)
const NERR_LogonsPaused = (NERR_BASE + 109)
const NERR_LogonServerConflict = (NERR_BASE + 110)
const NERR_LogonNoUserPath = (NERR_BASE + 111)
const NERR_LogonScriptError = (NERR_BASE + 112)
const NERR_StandaloneLogon = (NERR_BASE + 114)
const NERR_LogonServerNotFound = (NERR_BASE + 115)
const NERR_LogonDomainExists = (NERR_BASE + 116)
const NERR_NonValidatedLogon = (NERR_BASE + 117)
const NERR_ACFNotFound = (NERR_BASE + 119)
const NERR_GroupNotFound = (NERR_BASE + 120)
const NERR_UserNotFound = (NERR_BASE + 121)
const NERR_ResourceNotFound = (NERR_BASE + 122)
const NERR_GroupExists = (NERR_BASE + 123)
const NERR_UserExists = (NERR_BASE + 124)
const NERR_ResourceExists = (NERR_BASE + 125)
const NERR_NotPrimary = (NERR_BASE + 126)
const NERR_ACFNotLoaded = (NERR_BASE + 127)
const NERR_ACFNoRoom = (NERR_BASE + 128)
const NERR_ACFFileIOFail = (NERR_BASE + 129)
const NERR_ACFTooManyLists = (NERR_BASE + 130)
const NERR_UserLogon = (NERR_BASE + 131)
const NERR_ACFNoParent = (NERR_BASE + 132)
const NERR_CanNotGrowSegment = (NERR_BASE + 133)
const NERR_SpeGroupOp = (NERR_BASE + 134)
const NERR_NotInCache = (NERR_BASE + 135)
const NERR_UserInGroup = (NERR_BASE + 136)
const NERR_UserNotInGroup = (NERR_BASE + 137)
const NERR_AccountUndefined = (NERR_BASE + 138)
const NERR_AccountExpired = (NERR_BASE + 139)
const NERR_InvalidWorkstation = (NERR_BASE + 140)
const NERR_InvalidLogonHours = (NERR_BASE + 141)
const NERR_PasswordExpired = (NERR_BASE + 142)
const NERR_PasswordCantChange = (NERR_BASE + 143)
const NERR_PasswordHistConflict = (NERR_BASE + 144)
const NERR_PasswordTooShort = (NERR_BASE + 145)
const NERR_PasswordTooRecent = (NERR_BASE + 146)
const NERR_InvalidDatabase = (NERR_BASE + 147)
const NERR_DatabaseUpToDate = (NERR_BASE + 148)
const NERR_SyncRequired = (NERR_BASE + 149)
const NERR_UseNotFound = (NERR_BASE + 150)
const NERR_BadAsgType = (NERR_BASE + 151)
const NERR_DeviceIsShared = (NERR_BASE + 152)
const NERR_NoComputerName = (NERR_BASE + 170)
const NERR_MsgAlreadyStarted = (NERR_BASE + 171)
const NERR_MsgInitFailed = (NERR_BASE + 172)
const NERR_NameNotFound = (NERR_BASE + 173)
const NERR_AlreadyForwarded = (NERR_BASE + 174)
const NERR_AddForwarded = (NERR_BASE + 175)
const NERR_AlreadyExists = (NERR_BASE + 176)
const NERR_TooManyNames = (NERR_BASE + 177)
const NERR_DelComputerName = (NERR_BASE + 178)
const NERR_LocalForward = (NERR_BASE + 179)
const NERR_GrpMsgProcessor = (NERR_BASE + 180)
const NERR_PausedRemote = (NERR_BASE + 181)
const NERR_BadReceive = (NERR_BASE + 182)
const NERR_NameInUse = (NERR_BASE + 183)
const NERR_MsgNotStarted = (NERR_BASE + 184)
const NERR_NotLocalName = (NERR_BASE + 185)
const NERR_NoForwardName = (NERR_BASE + 186)
const NERR_RemoteFull = (NERR_BASE + 187)
const NERR_NameNotForwarded = (NERR_BASE + 188)
const NERR_TruncatedBroadcast = (NERR_BASE + 189)
const NERR_InvalidDevice = (NERR_BASE + 194)
const NERR_WriteFault = (NERR_BASE + 195)
const NERR_DuplicateName = (NERR_BASE + 197)
const NERR_DeleteLater = (NERR_BASE + 198)
const NERR_IncompleteDel = (NERR_BASE + 199)
const NERR_MultipleNets = (NERR_BASE + 200)
const NERR_NetNameNotFound = (NERR_BASE + 210)
const NERR_DeviceNotShared = (NERR_BASE + 211)
const NERR_ClientNameNotFound = (NERR_BASE + 212)
const NERR_FileIdNotFound = (NERR_BASE + 214)
const NERR_ExecFailure = (NERR_BASE + 215)
const NERR_TmpFile = (NERR_BASE + 216)
const NERR_TooMuchData = (NERR_BASE + 217)
const NERR_DeviceShareConflict = (NERR_BASE + 218)
const NERR_BrowserTableIncomplete = (NERR_BASE + 219)
const NERR_NotLocalDomain = (NERR_BASE + 220)
const NERR_IsDfsShare = (NERR_BASE + 221)
const NERR_DevInvalidOpCode = (NERR_BASE + 231)
const NERR_DevNotFound = (NERR_BASE + 232)
const NERR_DevNotOpen = (NERR_BASE + 233)
const NERR_BadQueueDevString = (NERR_BASE + 234)
const NERR_BadQueuePriority = (NERR_BASE + 235)
const NERR_NoCommDevs = (NERR_BASE + 237)
const NERR_QueueNotFound = (NERR_BASE + 238)
const NERR_BadDevString = (NERR_BASE + 240)
const NERR_BadDev = (NERR_BASE + 241)
const NERR_InUseBySpooler = (NERR_BASE + 242)
const NERR_CommDevInUse = (NERR_BASE + 243)
const NERR_InvalidComputer = (NERR_BASE + 251)
const NERR_MaxLenExceeded = (NERR_BASE + 254)
const NERR_BadComponent = (NERR_BASE + 256)
const NERR_CantType = (NERR_BASE + 257)
const NERR_TooManyEntries = (NERR_BASE + 262)
const NERR_ProfileFileTooBig = (NERR_BASE + 270)
const NERR_ProfileOffset = (NERR_BASE + 271)
const NERR_ProfileCleanup = (NERR_BASE + 272)
const NERR_ProfileUnknownCmd = (NERR_BASE + 273)
const NERR_ProfileLoadErr = (NERR_BASE + 274)
const NERR_ProfileSaveErr = (NERR_BASE + 275)
const NERR_LogOverflow = (NERR_BASE + 277)
const NERR_LogFileChanged = (NERR_BASE + 278)
const NERR_LogFileCorrupt = (NERR_BASE + 279)
const NERR_SourceIsDir = (NERR_BASE + 280)
const NERR_BadSource = (NERR_BASE + 281)
const NERR_BadDest = (NERR_BASE + 282)
const NERR_DifferentServers = (NERR_BASE + 283)
const NERR_RunSrvPaused = (NERR_BASE + 285)
const NERR_ErrCommRunSrv = (NERR_BASE + 289)
const NERR_ErrorExecingGhost = (NERR_BASE + 291)
const NERR_ShareNotFound = (NERR_BASE + 292)
const NERR_InvalidLana = (NERR_BASE + 300)
const NERR_OpenFiles = (NERR_BASE + 301)
const NERR_ActiveConns = (NERR_BASE + 302)
const NERR_BadPasswordCore = (NERR_BASE + 303)
const NERR_DevInUse = (NERR_BASE + 304)
const NERR_LocalDrive = (NERR_BASE + 305)
const NERR_AlertExists = (NERR_BASE + 330)
const NERR_TooManyAlerts = (NERR_BASE + 331)
const NERR_NoSuchAlert = (NERR_BASE + 332)
const NERR_BadRecipient = (NERR_BASE + 333)
const NERR_AcctLimitExceeded = (NERR_BASE + 334)
const NERR_InvalidLogSeek = (NERR_BASE + 340)
const NERR_BadUasConfig = (NERR_BASE + 350)
const NERR_InvalidUASOp = (NERR_BASE + 351)
const NERR_LastAdmin = (NERR_BASE + 352)
const NERR_DCNotFound = (NERR_BASE + 353)
const NERR_LogonTrackingError = (NERR_BASE + 354)
const NERR_NetlogonNotStarted = (NERR_BASE + 355)
const NERR_CanNotGrowUASFile = (NERR_BASE + 356)
const NERR_TimeDiffAtDC = (NERR_BASE + 357)
const NERR_PasswordMismatch = (NERR_BASE + 358)
const NERR_NoSuchServer = (NERR_BASE + 360)
const NERR_NoSuchSession = (NERR_BASE + 361)
const NERR_NoSuchConnection = (NERR_BASE + 362)
const NERR_TooManyServers = (NERR_BASE + 363)
const NERR_TooManySessions = (NERR_BASE + 364)
const NERR_TooManyConnections = (NERR_BASE + 365)
const NERR_TooManyFiles = (NERR_BASE + 366)
const NERR_NoAlternateServers = (NERR_BASE + 367)
const NERR_TryDownLevel = (NERR_BASE + 370)
const NERR_UPSDriverNotStarted = (NERR_BASE + 380)
const NERR_UPSInvalidConfig = (NERR_BASE + 381)
const NERR_UPSInvalidCommPort = (NERR_BASE + 382)
const NERR_UPSSignalAsserted = (NERR_BASE + 383)
const NERR_UPSShutdownFailed = (NERR_BASE + 384)
const NERR_BadDosRetCode = (NERR_BASE + 400)
const NERR_ProgNeedsExtraMem = (NERR_BASE + 401)
const NERR_BadDosFunction = (NERR_BASE + 402)
const NERR_RemoteBootFailed = (NERR_BASE + 403)
const NERR_BadFileCheckSum = (NERR_BASE + 404)
const NERR_NoRplBootSystem = (NERR_BASE + 405)
const NERR_RplLoadrNetBiosErr = (NERR_BASE + 406)
const NERR_RplLoadrDiskErr = (NERR_BASE + 407)
const NERR_ImageParamErr = (NERR_BASE + 408)
const NERR_TooManyImageParams = (NERR_BASE + 409)
const NERR_NonDosFloppyUsed = (NERR_BASE + 410)
const NERR_RplBootRestart = (NERR_BASE + 411)
const NERR_RplSrvrCallFailed = (NERR_BASE + 412)
const NERR_CantConnectRplSrvr = (NERR_BASE + 413)
const NERR_CantOpenImageFile = (NERR_BASE + 414)
const NERR_CallingRplSrvr = (NERR_BASE + 415)
const NERR_StartingRplBoot = (NERR_BASE + 416)
const NERR_RplBootServiceTerm = (NERR_BASE + 417)
const NERR_RplBootStartFailed = (NERR_BASE + 418)
const NERR_RPL_CONNECTED = (NERR_BASE + 419)
const NERR_BrowserConfiguredToNotRun = (NERR_BASE + 450)
const NERR_RplNoAdaptersStarted = (NERR_BASE + 510)
const NERR_RplBadRegistry = (NERR_BASE + 511)
const NERR_RplBadDatabase = (NERR_BASE + 512)
const NERR_RplRplfilesShare = (NERR_BASE + 513)
const NERR_RplNotRplServer = (NERR_BASE + 514)
const NERR_RplCannotEnum = (NERR_BASE + 515)
const NERR_RplWkstaInfoCorrupted = (NERR_BASE + 516)
const NERR_RplWkstaNotFound = (NERR_BASE + 517)
const NERR_RplWkstaNameUnavailable = (NERR_BASE + 518)
const NERR_RplProfileInfoCorrupted = (NERR_BASE + 519)
const NERR_RplProfileNotFound = (NERR_BASE + 520)
const NERR_RplProfileNameUnavailable = (NERR_BASE + 521)
const NERR_RplProfileNotEmpty = (NERR_BASE + 522)
const NERR_RplConfigInfoCorrupted = (NERR_BASE + 523)
const NERR_RplConfigNotFound = (NERR_BASE + 524)
const NERR_RplAdapterInfoCorrupted = (NERR_BASE + 525)
const NERR_RplInternal = (NERR_BASE + 526)
const NERR_RplVendorInfoCorrupted = (NERR_BASE + 527)
const NERR_RplBootInfoCorrupted = (NERR_BASE + 528)
const NERR_RplWkstaNeedsUserAcct = (NERR_BASE + 529)
const NERR_RplNeedsRPLUSERAcct = (NERR_BASE + 530)
const NERR_RplBootNotFound = (NERR_BASE + 531)
const NERR_RplIncompatibleProfile = (NERR_BASE + 532)
const NERR_RplAdapterNameUnavailable = (NERR_BASE + 533)
const NERR_RplConfigNotEmpty = (NERR_BASE + 534)
const NERR_RplBootInUse = (NERR_BASE + 535)
const NERR_RplBackupDatabase = (NERR_BASE + 536)
const NERR_RplAdapterNotFound = (NERR_BASE + 537)
const NERR_RplVendorNotFound = (NERR_BASE + 538)
const NERR_RplVendorNameUnavailable = (NERR_BASE + 539)
const NERR_RplBootNameUnavailable = (NERR_BASE + 540)
const NERR_RplConfigNameUnavailable = (NERR_BASE + 541)
const NERR_DfsInternalCorruption = (NERR_BASE + 560)
const NERR_DfsVolumeDataCorrupt = (NERR_BASE + 561)
const NERR_DfsNoSuchVolume = (NERR_BASE + 562)
const NERR_DfsVolumeAlreadyExists = (NERR_BASE + 563)
const NERR_DfsAlreadyShared = (NERR_BASE + 564)
const NERR_DfsNoSuchShare = (NERR_BASE + 565)
const NERR_DfsNotALeafVolume = (NERR_BASE + 566)
const NERR_DfsLeafVolume = (NERR_BASE + 567)
const NERR_DfsVolumeHasMultipleServers = (NERR_BASE + 568)
const NERR_DfsCantCreateJunctionPoint = (NERR_BASE + 569)
const NERR_DfsServerNotDfsAware = (NERR_BASE + 570)
const NERR_DfsBadRenamePath = (NERR_BASE + 571)
const NERR_DfsVolumeIsOffline = (NERR_BASE + 572)
const NERR_DfsNoSuchServer = (NERR_BASE + 573)
const NERR_DfsCyclicalName = (NERR_BASE + 574)
const NERR_DfsNotSupportedInServerDfs = (NERR_BASE + 575)
const NERR_DfsDuplicateService = (NERR_BASE + 576)
const NERR_DfsCantRemoveLastServerShare = (NERR_BASE + 577)
const NERR_DfsVolumeIsInterDfs = (NERR_BASE + 578)
const NERR_DfsInconsistent = (NERR_BASE + 579)
const NERR_DfsServerUpgraded = (NERR_BASE + 580)
const NERR_DfsDataIsIdentical = (NERR_BASE + 581)
const NERR_DfsCantRemoveDfsRoot = (NERR_BASE + 582)
const NERR_DfsChildOrParentInDfs = (NERR_BASE + 583)
const NERR_DfsInternalError = (NERR_BASE + 590)
const NERR_SetupAlreadyJoined = (NERR_BASE + 591)
const NERR_SetupNotJoined = (NERR_BASE + 592)
const NERR_SetupDomainController = (NERR_BASE + 593)
const NERR_DefaultJoinRequired = (NERR_BASE + 594)
const NERR_InvalidWorkgroupName = (NERR_BASE + 595)
const NERR_NameUsesIncompatibleCodePage = (NERR_BASE + 596)
const NERR_ComputerAccountNotFound = (NERR_BASE + 597)
const NERR_PersonalSku = (NERR_BASE + 598)
const NERR_PasswordMustChange = (NERR_BASE + 601)
const NERR_AccountLockedOut = (NERR_BASE + 602)
const NERR_PasswordTooLong = (NERR_BASE + 603)
const NERR_PasswordNotComplexEnough = (NERR_BASE + 604)
const NERR_PasswordFilterError = (NERR_BASE + 605)
const MAX_NERR = (NERR_BASE + 899)
