# Golang 单测最佳实践
## Golang 单测简介
Golang 内置 `testing` 包和 `go test` 命令，无需额外配置即可编写和运行测试。测试文件以 `_test.go` 结尾，测试函数以 Test 前缀命名（如 `TestAdd`），并通过 `*testing.T` 参数管理测试状态；也能使用 `t.Parallel()`和 `Bench` 前缀进行并发和性能测试。

Golang 单元测试在简洁性、高效性和丰富的工具链方面具备显著优势，并且便于与其他工具或系统进行集成，例如持续集成（CI）系统。

### 测试框架
+ [testing](https://pkg.go.dev/testing) 包
+ [testify](https://github.com/stretchr/testify/tree/master/suite) suite
+ [ginkgo](https://github.com/onsi/ginkgo)

### 常用对比工具
+ reflect.DeepEqual
+ [testify](https://github.com/stretchr/testify) required/assert
+ [gomega](https://github.com/onsi/gomega)

## 单测场景
### 1、表格驱动测试
表格驱动测试(Table-Driven Testing)将测试数据与测试逻辑分离，通常用结构体定义输入参数、预期输出和测试条件，逐条对测试用例进行验证，使测试代码更清晰，也能很方便的添加测试场景，更容易达到高测试覆盖率。

```go
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
```

```go
func TestFormatSize(t *testing.T) {
	cases := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "zero",
			size:     0,
			expected: "0 B",
		},
		{
			name:     "1 byte",
			size:     1,
			expected: "1 B",
		},
        ...
		{
			name:     "1000000000000000000 bytes",
			size:     1000000000000000000,
			expected: "888.2 PB",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			formatted := FormatSize(tc.size)
			if formatted != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, formatted)
			}
		})
	}
}
```

### 2、并行测试
在Go语言单元测试中，`t.Parallel()`允许测试用例并行执行。通过并行执行独立测试用例，充分利用多核CPU资源，能够提高测试效率，适用在比较耗时的测试场景，但是需要注意共享资源冲突问题。

```go
func TestFormatSize(t *testing.T) {
	cases := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "zero",
			size:     0,
			expected: "0 B",
		},
		{
			name:     "1 byte",
			size:     1,
			expected: "1 B",
		},
        ...
		{
			name:     "1000000000000000000 bytes",
			size:     1000000000000000000,
			expected: "888.2 PB",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
			formatted := FormatSize(tc.size)
			if formatted != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, formatted)
			}
		})
	}
}
```

```go
    for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
            tc := tc
            t.Parallel()
			formatted := FormatSize(tc.size)
			if formatted != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, formatted)
			}
		})
	}
```

`t.Parallel()`在 for-loop 中使用时要注意 golang loopvar 的陷阱，或者升级到最新的 golang 版本 Go 1.22+。

[https://go.dev/wiki/LoopvarExperiment](https://go.dev/wiki/LoopvarExperiment)

### 3、子测试（Subtests）
子测试 Subtest（包含基准测试）是在单个测试函数中创建嵌套的测试层次结构来定义子测试和子基准测试，而不必为每个子测试和子基准测试定义单独的函数，每个子测试和子基准测试都有一个唯一的名称。

```go
func TestFormatSize_subtests(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		formatted := FormatSize(0)
		if formatted != "0 B" {
			t.Errorf("expected \"0 B\", got \"%s\"", formatted)
		}
	})
	t.Run("100000", func(t *testing.T) {
		formatted := FormatSize(100000)
		if formatted != "97.7 KB" {
			t.Errorf("expected \"97.7 KB\", got \"%s\"", formatted)
		}
	})
}
```

### 4、初始化数据和清理
在Go语言单元测试中，经常会需要初始化一些全局数据，包括一些数据库连接、服务 Config 配置等等，为了确保测试可靠性和隔离性的同时能够复用资源，需要合理的初始化数据和清理。

#### 一、使用 TestMain 进行全局初始化
在 Golang 单测中可以编写 `func TestMain(m *testing.M)`方法（非必须）是作为测试代码的主入口，提供额外的 Setup/Teardown 操作。

```go
var (
	testUserRepo UserRepository
	sqlMock      sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
	mockDb, mock, _ := sqlmock.New()
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      mockDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("test open mock db failed: %v", err)
	}
	defer mockDb.Close()
	testUserRepo = NewUserRepository(db)
	sqlMock = mock

	m.Run()
}

```

#### 二、[testify suite](https://pkg.go.dev/github.com/stretchr/testify/suite) 的 Setup/Teardown
testify suite 包提供了测试的 Setup/Teardown 方法，SetupSuite/TearDownSuite 在套件全局生效，SetupTest/TearDownTest 对每个 Test 测试生效。

```go
type ExampleTestSuite struct {
    suite.Suite
    VariableThatShouldStartAtFive int
}

func (suite *ExampleTestSuite) SetupSuite() {
    suite.VariableThatShouldStartAtFive = 5
}

func (suite *ExampleTestSuite) SetupTest() {
    suite.VariableThatShouldStartAtFive = 3
}

func (suite *ExampleTestSuite) TearDownTest() {
}

func (suite *ExampleTestSuite) TearDownSuite() {
}
```

#### 三、测试方法内部 Setup
```go
func TestUserWithMock(t *testing.T) {
    mockDb, mock, _ := sqlmock.New()
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      mockDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("test open mock db failed: %v", err)
	}
	defer mockDb.Close()

    ...
}
```

### 5、帮助函数（helpers）
`t.Helper()` 是一个用于优化测试错误报告的方法，主要作用是将辅助函数中的错误定位到调用方（即实际测试用例）的位置。好处是减少重复的测试逻辑，通过封装辅助函数复用代码，同时保持错误信息的清晰性。

```go
func TestExample(t *testing.T) {
    assertEqual(t, 1, 2)
}

func assertEqual(t *testing.T, a, b int) {
    t.Helper()

    if a != b {
        t.Errorf("assert failed: %d != %d", a, b)
    }

    assertErr(t, nil)
}

func assertErr(t testing.T, err error){
    t.Helper()

    t.Error(err)
}
```

```go
func assertConfig(t *testing.T, expectedConfig *Config, conf *Config) {
	t.Helper()

	if expectedConfig == nil || conf == nil {
		if expectedConfig != nil || conf != nil {
			t.Errorf("assert config not equal, expected: %v, actual: %v", expectedConfig, conf)
		}
		return
	}
	if expectedConfig.DBHost != conf.DBHost ||
		expectedConfig.DBPort != conf.DBPort ||
		expectedConfig.DBUser != conf.DBUser ||
		expectedConfig.DBPassword != conf.DBPassword ||
		expectedConfig.DBName != conf.DBName ||
		expectedConfig.ListenPort != conf.ListenPort ||
		expectedConfig.PprofAddr != conf.PprofAddr {
		t.Errorf("assert config not equal, expected: %v, actual: %v", expectedConfig, conf)
	}
}
```

### 6、Skip 跳过
`t.Skip` 可以在测试函数中调用，用于跳过当前测试用例。通常用于在某些条件不满足时跳过测试，例如环境变量未设置或依赖服务不可用。

```go
func TestFormatSize_subtests(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		formatted := FormatSize(0)
		if formatted != "0 B" {
			t.Errorf("expected \"0 B\", got \"%s\"", formatted)
		}
	})
	t.Run("100000", func(t *testing.T) {
		formatted := FormatSize(100000)
		if formatted != "97.7 KB" {
			t.Errorf("expected \"97.7 KB\", got \"%s\"", formatted)
		}
	})
	t.Run("-100", func(t *testing.T) {
		t.Skip()
	})
}
```

如果测试依赖某些环境或者条件，在检查依赖条件不足时使用 `t.Skip`跳过此测试。

```go
func (s *UserTestSuite) SetupSuite() {
	dbhost := os.Getenv("TEST_DBHOST")
	if dbhost == "" {
		s.T().Skip("skip test: env TEST_DBHOST not set")
	}
	dbportStr := os.Getenv("TEST_DBPORT")
	dbport := 3306
	if dbportStr != "" {
		var err error
		dbport, err = strconv.Atoi(dbportStr)
		if err != nil {
			s.T().Skipf("skip test: parse dbport from ENV failed: %s", dbportStr)
		}
	}
	dbuser := os.Getenv("TEST_DBUSER")
	dbpassword := os.Getenv("TEST_DBPASSWORD")
	dbname := os.Getenv("TEST_DBNAME")
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		dbuser, dbpassword, dbhost, dbport, dbname)), &gorm.Config{})
	if err != nil {
		s.T().Skipf("skip test: open database failed: %v", err)
	}
	s.userRepo = NewUserRepository(db)
}
```

### 7、临时目录/文件
在我们代码编写中经常会遇到一些文件读写操作，单元测试时为了隔离测试环境、避免污染生产数据或解决路径依赖问题，我们创建临时文件或者临时目录，结束后进行清理。

#### 一、使用标准库 `os`
通过 `os.MkdirTemp` 和 `os.CreateTemp` 生成临时目录和文件，测试结束后需手动清理。

#### 二、使用 testing t.TempDir()
Go 1.15+ 后在 testing 包新增了 `TempDir` 方法，并且在测试函数运行结束后自动清理。

```go
func TestLoadConfig(t *testing.T) {
	// dir, err := os.MkdirTemp("", "testloadconfig")
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }
	// defer os.RemoveAll(dir)

	// after go1.15
	dir := t.TempDir()
    ...
}
```

### 8、基准测试（Benchmark）
基准测试是 golang 中对代码性能（内存分配、执行时间、并发性）进行评估的测试工具， 并且能够根据测试结果不断优化以提升性能。基准测试函数形式 `func BenchmarkXxx(*testing.B)`，在运行 `go test`时使用 `-bench`运行。

```go
func BenchmarkFormatSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatSize(100000000000)
	}
}
```

#### ResetTimer 和 RunParallel
为了避免初始化数据对基准结果有影响，需要使用 ResetTimer 重新计时和内存申请统计。

使用 RunParallel 能够利用多核特性并行运行，验证多 goroutines（默认取决于GOMAXPROCS）下的并发性能。

```go
func BenchmarkFormatSize(b *testing.B) {
    // initialize some data
    // b.ResetTimer before run benchmark
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            FormatSize(100000000000)
        }
    })
}
```

详细的 test flag 参考：[https://pkg.go.dev/cmd/go#hdr-Testing_flags](https://pkg.go.dev/cmd/go#hdr-Testing_flags)

```bash
$ go test -bench=. -benchtime=5s -benchmem
```

### 9、HTTP 测试
Golang 提供了 [iotest](https://pkg.go.dev/testing/iotest) 和 [httptest](https://pkg.go.dev/net/http/httptest) 包，方便测试网络 IO 和 HTTP 服务。提供了常用的

HTTP handler 和 Server 测试方法。

#### HTTP Handler
测试 handler `func(w ResponseWriter, r *Request)`需要 mock `ResponseWriter` 和 `*Request` 实例，使用 `httptest.NewRequest` 和`httptest.NewRecorder()`进行初始化，在测试代码运行后可以对 Recorder 数据进行验证。

```go
func (s *ServiceTestSuite) TestCreateUser() {
	s.Run("param name empty", func() {
		req := httptest.NewRequest("POST", "http://127.0.0.1:8888/user/create?email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusBadRequest, w.Code)
		s.EqualValues(`{"error":"param name not set"}`, w.Body.String())
	})
	s.Run("success", func() {
		t := time.Unix(1752999201, 0)
		id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
		s.mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
		s.mockUserRepo.EXPECT().GetByEmail(gomock.Any()).Return(&store.User{
			ID:        id,
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: t,
			UpdatedAt: t,
		}, nil).Times(1)

		req := httptest.NewRequest("POST", "http://127.0.0.1:8888/user/create?name=liuliu&email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusOK, w.Code)
		s.EqualValues(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`, w.Body.String())
	})
}
```

#### HTTP Client
Client 代码中需要 http 请求，我们使用 `httptest.NewServer`创建测试需要的临时 http server，将 server.URL 传给 client 进行请求，并在测试完成后使用 `server.Close()`关闭服务。

在测试过程中对 server 添加自定义 handler 对 server 请求处理进行 mock。

```go
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r)
	}))
	defer server.Close()

	c := New(server.URL)

	t.Run("create", func(t *testing.T) {
		handleFunc = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`))
		}
		user, err := c.UserCreate(User{
			Name:  "liuliu",
			Email: "aa@bb.com",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, &User{
			ID:        "0198271f-bc9d-74ac-a63b-41cf2c6c2f82",
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: time.Unix(1752999201, 0),
			UpdatedAt: time.Unix(1752999201, 0),
		}, user)
	})
```

### 10、连接数据库
#### 使用 [sqlmock](https://github.com/DATA-DOG/go-sqlmock)
**sqlmock** 是一个实现 [sql/driver](https://godoc.org/database/sql/driver) 的模拟库。它能模拟测试中的任何 **sql** 驱动程序行为，而无需真正的数据库连接。

缺点是 mock 数据繁琐，如果是上层业务代码更适合用 gomock，如果是验证 orm 代码更适合用 integration-testing。

```go
mockDb, mock, _ := sqlmock.New()
dialector := mysql.New(mysql.Config{
    DSN:                       "sqlmock_db_0",
    DriverName:                "mysql",
    Conn:                      mockDb,
    SkipInitializeWithVersion: true,
})
db, err := gorm.Open(dialector, &gorm.Config{})
```

```go
	t.Run("Create", func(t *testing.T) {
		u := &User{
			ID:    uuid.NewString(),
			Name:  "liuhong",
			Email: "aaa@bb.com",
			Age:   22,
		}
		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()
		err := testUserRepo.Create(u)
		assert.NoError(t, err)
	})
```

#### 使用 sqlite
如果你使用的 sql 没有数据库专用函数或语句，可以使用 sqlite 更加方便快捷。

```go
import (
	"gorm.io/driver/sqlite"
	_ "github.com/glebarez/go-sqlite"
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    ...
}
```

#### integration 方式
使用集成测试连接真实的数据库，在验证数据库相关操作的时候更加方便和准确。

```go
func (s *UserTestSuite) SetupSuite() {
	dbhost := os.Getenv("TEST_DBHOST")
	if dbhost == "" {
		s.T().Skip("skip test: env TEST_DBHOST not set")
	}
	dbportStr := os.Getenv("TEST_DBPORT")
	dbport := 3306
	if dbportStr != "" {
		var err error
		dbport, err = strconv.Atoi(dbportStr)
		if err != nil {
			s.T().Skipf("skip test: parse dbport from ENV failed: %s", dbportStr)
		}
	}
	dbuser := os.Getenv("TEST_DBUSER")
	dbpassword := os.Getenv("TEST_DBPASSWORD")
	s.dbname = fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	// create test database
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8mb4&parseTime=true&loc=Local",
		dbuser, dbpassword, dbhost, dbport)), &gorm.Config{})
	if err != nil {
		s.T().Skipf("skip test: open database failed: %v", err)
	}
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", s.dbname)).Error
	if err != nil {
		s.T().Skipf("skip test: create database failed: %v", err)
	}

	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbuser, dbpassword, dbhost, dbport, s.dbname)), &gorm.Config{})
	if err != nil {
		s.T().Skipf("skip test: open database failed: %v", err)
	}
	s.db = db
	s.userRepo = NewUserRepository(db)

	s.db.AutoMigrate(&User{})
}

func (s *UserTestSuite) TearDownSuite() {
	s.db.Exec(fmt.Sprintf("DROP DATABASE %s", s.dbname))
}
```

### 11、Mock 接口
#### [gomock](https://github.com/uber-go/mock)
gomock 是一套 golang mock 测试框架，需要使用 mockgen 工具对 interface 生成 mock 代码，而后在测试代码使用 mock 代码。

生成 mock 代码：`go generate ./...`。

```go
//go:generate mockgen -source=user.go -destination=user_mock.go -package=store
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	DeleteByID(id string) error
	List(page, pageSize int) ([]User, int64, error)
}
```

调用 mock 实例的 `EXPECT()`方法返回 mock recorder，使用 recorder 对期望调用的方法、参数、返回数据、调用次数等等进行设置。

```go
func (s *ServiceTestSuite) TestCreateUser() {
	s.Run("param name empty", func() {
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/user/create?email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusBadRequest, w.Code)
		s.EqualValues(`{"error":"param name not set"}`, w.Body.String())
	})
	s.Run("success", func() {
		t := time.Unix(1752999201, 0)
		id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
		s.mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
		s.mockUserRepo.EXPECT().GetByEmail(gomock.Any()).Return(&store.User{
			ID:        id,
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: t,
			UpdatedAt: t,
		}, nil).Times(1)

		req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/user/create?name=liuliu&email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusOK, w.Code)
		s.EqualValues(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`, w.Body.String())
	})
}
```

#### fake/test 抽象实现
在 mock 不方便或者调用繁琐的情况下，为了测试我们也可以做 fake 实现，在调用多次的情况下能够省略大量 mock 数据生成。

```go
//go:generate mockgen -source=client.go -destination=client_mock.go -package=client
type Client interface {
	UserCreate(u User) (*User, error)
	UserGet(id string) (*User, error)
	UserUpdate(u User) error
	UserDelete(id string) error
	UserList() ([]User, int64, error)
}
```

```go
type fakeClient struct {
	mu           sync.Mutex
	users        map[string]*client.User
	usersByEmail map[string]string
}

var _ client.Client = &fakeClient{}

func (c *fakeClient) UserCreate(u client.User) (*client.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.usersByEmail[u.Email]
	if ok {
		return nil, fmt.Errorf("user already exists")
	}

	id := uuid.NewString()
	now := time.Now()
	c.users[id] = &client.User{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}
	c.usersByEmail[u.Email] = id
	return c.users[id], nil
}
```

例如为了测试 `io.Reader` 中返回错误，实现 testIOReader ，自定义返回错误（简单的也可以使用 [iotest](https://pkg.go.dev/testing/iotest) 包），比如在返回部分数据后返回 `io.ErrUnexpectedEOF`。

```go
func testIOReader struct{
    data []byte
    err error
}

func (r *testIOReader)Read(p []byte) (n int, err error){
    ...
    return r.err
}
```

#### [gomonkey](https://github.com/agiledragon/gomonkey)
gomonkey 是一个对 golang 方法/函数进行 patch 的测试库，它能够在测试时实时替换/还原方法和函数并返回期望的数据，以达到 mock 的效果，但是又不需要提前生成 mock 等测试代码。

因为 CPU 和编译相关原因，运行需要添加 `-gcflags=all=-l`参数避免失败。

原理查看：[https://bou.ke/blog/monkey-patching-in-go/](https://bou.ke/blog/monkey-patching-in-go/)

```go
	t.Run("list", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodReturn(&http.Client{}, "Do", &http.Response{
			StatusCode: http.StatusOK,
		}, nil)
		patches.ApplyFuncReturn(io.ReadAll, []byte(`{"data":{"total":1,"users":[{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}]}}`), nil)

		users, total, err := c.UserList()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, total)
		assert.EqualValues(t, []User{{
			ID:        "0198271f-bc9d-74ac-a63b-41cf2c6c2f82",
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: time.Unix(1752999201, 0),
			UpdatedAt: time.Unix(1752999201, 0),
		}}, users)
	})
```

## 单测覆盖率
测试覆盖率是衡量代码质量的一个重要指标，提高覆盖率能够显著降低生产环境中的bug发生率。但是覆盖率也不是越高越好，在覆盖率达到一定比例之后提高覆盖率所需要的成本也会越来越高，Golang 因为语言特性 if-err 判断比较多也会影响覆盖率，建议不同功能级别使用不同的覆盖率。

核心库：90%； 核心逻辑代码：80%；常规业务代码：60%。

在编写测试代码中避免**虚假覆盖率**（执行但不验证的测试） 和过度追求 100% 覆盖率。

### 查看覆盖率
#### 查看包覆盖率
```bash
$ go test -cover
PASS
coverage: 63.7% of statements
ok      go-unittest-best-practice/internal/api 0.557s
```

#### 查看项目整体覆盖率
```bash
# 生成 coverage.out
$ go test -coverprofile=cover.out ./...
# 查看全部方法覆盖率
$ go tool cover -func=cover.out
go-unittest-best-practice/cmd/user_manage/main.go:22:          main                    0.0%
go-unittest-best-practice/internal/api/service.go:22:          NewService              100.0%
...
go-unittest-best-practice/pkg/client/fake/fake_client.go:79:   UserList                0.0%
go-unittest-best-practice/pkg/utils/format.go:5:               FormatSize              100.0%
total:                                                  (statements)            39.3%
```

total 对应的是项目整体覆盖率，但是把 mock 文件也算进了正常代码文件，导致覆盖率偏低。

#### 过滤不需要覆盖的代码
```bash
$ cat cover.out | grep -v -E "mock.go|fake|main.go" > cover2.out
$ go tool cover -func=cover2.out
go-unittest-best-practice/internal/api/service.go:22:          NewService              100.0%
go-unittest-best-practice/internal/api/service.go:37:          ServeHTTP               100.0%
...
go-unittest-best-practice/pkg/client/client.go:172:            UserList                57.1%
go-unittest-best-practice/pkg/utils/format.go:5:               FormatSize              100.0%
total:                                                  (statements)            61.6%
```

#### 查看行覆盖
```bash
# 浏览器打开覆盖行统计 coverage.html
$ go tool cover -html cover2.out
# 或者手动生成 coverage.html
$ go tool cover -html cover2.out -o coverage.html
```

### 提高覆盖率
在代码结构优化上可以考虑几点以提高代码的可测试性。

#### 使用 Table-Driven Testing
表格驱动测试(Table-Driven Testing) 便于添加新的测试用例而不需要修改测试逻辑，能够快速测试场景，有助于提高单测覆盖率。

#### 复杂逻辑拆分
复杂的逻辑代码不只不利于阅读，在编写单测用例上也比较困难，将复杂函数拆分为多个小函数，每个函数职责单一，便于单独测试。

#### 抽象接口
如果你实现的代码需要作为依赖被调用，可以考虑定义抽象接口，让上层调用依赖接口而非具体实现，方便上层对依赖进行 mock。

#### 依赖注入
依赖注入将对象的依赖关系从内部创建改为外部注入，例如使用参数传递，减少了组件间的直接耦合，可以轻松注入Mock对象进行单元测试提高可测试性。

## 工具使用
```bash
# 测试单个用例，-v 表示Verbose output，并且使用 -run 根据前缀匹配
$ go test -v -run TestService

# 跳过用例
$ go test -v -skip TestService/TestListUser

# 生成 coverage.out
$ go test -coverprofile=cover.out

# 运行整个项目的测试用例
$ go test ./...

# 生成整个项目的测试用例的 coverage.out
$ go test -coverprofile=cover.out ./...

# 查看 func 覆盖率
$ go tool cover -func=cover.out

# 浏览器打开覆盖行统计 coverage.html
$ go tool cover -html cover.out

# 手动生成 coverage.html
$ go tool cover -html cover.out -o coverage.html

# 运行 benchmark，-bench也支持前缀匹配
$ go test -bench=.

# benchmark 生成 mem 相关统计
$ go test -bench=. -benchmem

# benchmark 运行时长，默认1s
$ go test -bench=. -benchmem -benchtime 5s

# benchmark 时生成 cpu/mem profile
$ go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out

# 更加方便的运行测试的工具 gotestsum
# https://github.com/gotestyourself/gotestsum
$ gotestsum --format testname
```

## 参考
- [https://go.dev/doc/tutorial/add-a-test](https://go.dev/doc/tutorial/add-a-test)
- [https://pkg.go.dev/cmd/go#hdr-Testing_flags](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [https://pkg.go.dev/testing](https://pkg.go.dev/testing)
- [https://pkg.go.dev/testing/iotest](https://pkg.go.dev/testing/iotest)
- [https://pkg.go.dev/net/http/httptest](https://pkg.go.dev/net/http/httptest)
- [https://github.com/uber-go/mock](https://github.com/uber-go/mock)
- [https://github.com/agiledragon/gomonkey](https://github.com/agiledragon/gomonkey)
- [https://github.com/prashantv/gostub](https://github.com/prashantv/gostub)
- [https://github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
