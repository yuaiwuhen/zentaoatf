package commConsts

type ExecFromDef string

const (
	FromZentao ExecFromDef = "zentao"
	FromCmd    ExecFromDef = "cmd"
	FromClient ExecFromDef = "client"
)

func (e ExecFromDef) String() string {
	return string(e)
}

type ResponseCode struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

var (
	Success         = ResponseCode{0, "Request Successfully"}
	CommErr         = ResponseCode{100, "Common Error"}
	ParamErr        = ResponseCode{200, "Parameter Error"}
	UnAuthorizedErr = ResponseCode{401, "UnAuthorized"}

	ResponseParseErr = ResponseCode{6000, "Json Parse Error"}
	NeedInitErr      = ResponseCode{1000, "Data Not Init"}

	ErrDirNotEmpty                    = ResponseCode{1100, "Dir Not Empty Error"}
	ErrZentaoConfig                   = ResponseCode{2000, "Zentao Config Error"}
	ErrZentaoRequest                  = ResponseCode{3000, "zentao request Error"}
	ErrRecordWithSameNameExist        = ResponseCode{3010, "Record With Same Name Exist"}
	ErrRecordWithSamePathExist        = ResponseCode{3011, "Record With Same Path Exist"}
	ErrRecordWithSameInterpreterExist = ResponseCode{2000, "Record With Same Interpreter Exist"}
	ErrRecordNotExist                 = ResponseCode{4000, "Record Not Found"}
	ErrRecordNotAllowDelete           = ResponseCode{2000, "Record Can Not Delete"}
	ErrPostParam                      = ResponseCode{2000, "Parameter Error"}
	NotAvailable                      = ResponseCode{5000, "Not Available"}
	ServerNotAvailable                = ResponseCode{5001, "Server Not Available"}
	ProxyNotAvailable                 = ResponseCode{5002, "Proxy Not Available"}

	ErrFileOrDirExist = ResponseCode{6000, "File or Dir Exist"}
)

type ResultStatus string

const (
	PASS    ResultStatus = "pass"
	FAIL    ResultStatus = "fail"
	SKIP    ResultStatus = "skip"
	BLOCKED ResultStatus = "blocked"
	UNKNOWN ResultStatus = "unknown"
)

func (e ResultStatus) String() string {
	return string(e)
}

type ExecCmd string

const (
	ExecInit ExecCmd = "init"
	ExecStop ExecCmd = "execStop"

	ExecCase   ExecCmd = "execCase"
	ExecModule ExecCmd = "execModule"
	ExecSuite  ExecCmd = "execSuite"
	ExecTask   ExecCmd = "execTask"

	ExecUnit    ExecCmd = "execUnit"
	ExecZapScan ExecCmd = "zapScan"

	Watch ExecCmd = "watch"
)

func (e ExecCmd) String() string {
	return string(e)
}

type WsMsgCategory string

const (
	Output WsMsgCategory = "output"

	Run    WsMsgCategory = "run"
	Result WsMsgCategory = "result"
	Error  WsMsgCategory = "error"

	Communication WsMsgCategory = "communication"
	Unknown       WsMsgCategory = ""
)

func (e WsMsgCategory) String() string {
	return string(e)
}

type OsType string

const (
	OsWindows OsType = "windows"
	OsLinux   OsType = "linux"
	OsMac     OsType = "mac"
)

func (e OsType) String() string {
	return string(e)
}
func (OsType) Get(osName string) OsType {
	return OsType(osName)
}

type CaseStepType string

const (
	Group CaseStepType = "group"
	Item  CaseStepType = "item"
)

type ExecBy string

const (
	Case   ExecBy = "case"
	Module ExecBy = "module"
	Suite  ExecBy = "suite"
	Task   ExecBy = "task"
)

func (e ExecBy) String() string {
	return string(e)
}
func (ExecBy) Get(str string) ExecBy {
	return ExecBy(str)
}

type TestType string

const (
	TestFunc TestType = "func"
	TestUnit TestType = "unit"
	TestZap  TestType = "zap"
)

func (e TestType) String() string {
	return string(e)
}
func (TestType) Get(str string) TestType {
	return TestType(str)
}

type TestTool string

const (
	ZTF     TestTool = "ztf"
	Allure  TestTool = "allure"
	JUnit   TestTool = "junit"
	TestNG  TestTool = "testng"
	PHPUnit TestTool = "phpunit"
	PyTest  TestTool = "pytest"
	Jest    TestTool = "jest"
	CppUnit TestTool = "cppunit"
	GTest   TestTool = "gtest"
	QTest   TestTool = "qtest"
	GoTest  TestTool = "gotest"

	AutoIt         TestTool = "autoit"
	Selenium       TestTool = "selenium"
	Appium         TestTool = "appium"
	RobotFramework TestTool = "robot"
	Cypress        TestTool = "cypress"
	Playwright     TestTool = "playwright"
	Puppeteer      TestTool = "puppeteer"
	K6             TestTool = "k6"
	Zap            TestTool = "zap"
)

func (e TestTool) String() string {
	return string(e)
}
func (TestTool) Get(str string) TestTool {
	return TestTool(str)
}

type BuildTool string

const (
	Maven BuildTool = "maven"
	Mocha BuildTool = "mocha"
)

func (e BuildTool) String() string {
	return string(e)
}
func (BuildTool) Get(str string) BuildTool {
	return BuildTool(str)
}

type PlatformType string

const (
	Android PlatformType = "android"
	Ios     PlatformType = "ios"
	Host    PlatformType = "host"
	Vm      PlatformType = "vm"
)

type TreeNodeType string

const (
	Root      TreeNodeType = "root"
	Workspace TreeNodeType = "workspace"
	Dir       TreeNodeType = "dir"
	File      TreeNodeType = "file"

	ZentaoModule TreeNodeType = "module"
)

type ScriptFilterType string

const (
	FilterWorkspace ScriptFilterType = "workspace"
	FilterModule    ScriptFilterType = "module"
	FilterSuite     ScriptFilterType = "suite"
	FilterTask      ScriptFilterType = "task"
)

type NodeCreateMode string

const (
	Brother NodeCreateMode = "brother"
	Child   NodeCreateMode = "child"
)

func (e NodeCreateMode) String() string {
	return string(e)
}

type NodeCreateType string

const (
	CreateDir  NodeCreateType = "dir"
	CreateNode NodeCreateType = "node"
)

func (e NodeCreateType) String() string {
	return string(e)
}

type DropPos int

const (
	Before DropPos = -1
	Inner  DropPos = 0
	After  DropPos = 1
)

func (e DropPos) Int() int {
	return int(e)
}

type JobStatus string

const (
	JobCreated    JobStatus = "created"
	JobInprogress JobStatus = "inprogress"

	// temp status
	JobTimeout JobStatus = "timeout"
	JobError   JobStatus = "error"

	// final status
	JobCanceled  JobStatus = "canceled"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

func (e JobStatus) ToString() string {
	return string(e)
}

type K6ResultType string

const (
	Metric K6ResultType = "Metric"
	Point  K6ResultType = "Point"
)

func (e K6ResultType) ToString() string {
	return string(e)
}
