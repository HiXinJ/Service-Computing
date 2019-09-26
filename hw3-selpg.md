# Golang实现CLI实用程序selpg

## 项目连接

[goOnline](http://139.9.57.167:20080/share/bm67l1ud0liuiksqatg0?secret=false)  
[github](https://github.com/HiXinJ/Service-Computing/tree/master/selpg)

## selpg简介
selpg名称代表Select Page。该实用程序从标准输入或从命令行参数给出的文件名读取文本输入。它允许用户指定来自该输入并随后将被输出的页面范围。例如，如果输入含有 100 页，则用户可指定只打印第 35 至 65 页。这种特性有实际价值，因为在打印机上打印选定的页面避免了浪费纸张。另一个示例是，原始文件很大而且以前已打印过，但某些页面由于打印机卡住或其它原因而没有被正确打印。在这样的情况下，则可用该工具来只打印需要打印的页面。

## 命令参数 

- “-s Number”和“-e Number”强制选项：  
selpg 要求用户用两个命令行参数“-s Number”（例如，“-s10”表示从第 10 页开始）和“-e Number”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。selpg 对所给的页号进行合理性检查；换句话说，它会检查两个数字是否为有效的正整数以及结束页是否不小于起始页。这两个选项，“-sNumber”和“-eNumber”是强制性的. 
例如：
    ```bash
    $ selpg -s1 -e10 infile
    ```
    将infile第1-10页输出到屏幕

- “-lNumber”和“-f”可选选项：  
selpg 可以处理两种输入文本：  
**类型1**：该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-l Number”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。该缺省值可以用“-lNumber”选项覆盖，如下所示： 

    ```shell
    $ selpg -s1 -e10 -l 66 infile
    ```
    **类型2**: 格式由“-f”选项表示，如下所示：
    ```
    $ selpg -s10 -e20 -f infile
    ```
    该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。  

    ***注：“-lNumber”和“-f”选项是互斥的。***

- “-dDestination”可选选项：  
selpg 还允许用户使用“-dDestination”选项将选定的页直接发送至打印机。这里，“Destination”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。该目的地应该存在 ― selpg 不检查这一点。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“Destination”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出。这一特性是用exec.Command()实现的，在下面的示例中，我们打开到命令

    ```
    $ lp -dDestination
    ```
    的管道以便输出，并写至该管道而不是标准输出：

    ```
    selpg -s10 -e20 -dlp1
    ```

## selpg的实现  

下面介绍实现selpg命令行程序用到的Golang包  
### 1、pflag
pflag用于解析输入命令参数  
下面代码定义一个Int型flag，参数为flag名称，简写，默认值，用法。返回为Int型指针。  
Parse用于对命令行输入的参数解析，即解析os.Args[1:]
```go
var pstartPage = pflag.IntP("startPage", "s", -1, "start page")
pflag.Parse()
```

### 2、bufio
bufio用来读取文件  
```go
if fin, err := os.Open(sa.inFilename); err != nil {
			fmt.Println(err)
			os.Exit(1)
}else{
    finReader := bufio.NewReader(fin)    
    line, err = finReader.ReadBytes('\n')//逐行读取
}
```

### 3、os/exec
exec.command()用于实现lp命令的管道输入。  
下面代码将selpg的输出定位到lp命令的输入：
```go
cmd = exec.Command("lp", "-d", sa.printDest)
cmd.Stderr = os.Stderr
fout, _ = cmd.StdinPipe()
```
开始命令并等待其结束：
```go
cmd.Run()
```
### 4、命令参数的数据结构  
```go
type selpgArgs struct {
	startPage  int
	endPage    int
	inFilename string
	pageLen    int
	pageType   bool //True for form-feed-delimited page type, fasle for lines-delimited page type
	printDest  string
}
```
### 5、两个主要函数
ValidateArgs会对传入的参数以及选项sa进行验证，如要满足： 
- 参数<=1, 只有inFilename参数从文件输入；或0参数，即标准输入
- startPage >= 1 
- startPage >= endPage
- pageLen >= 1
```go
func ValidateArgs(sa selpgArgs) {
	if pflag.NArg() > 1 {
		pflag.Usage()
		os.Exit(1)
	}

	if sa.startPage < 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", progname, sa.startPage)
		os.Exit(1)
	}
	if sa.startPage > sa.endPage {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", progname, sa.endPage)
		os.Exit(1)
	}
	if sa.pageLen <= 0 {
		fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", progname, sa.pageLen)
		os.Exit(1)
	}
}
```

processInput获得命令参数信息sa，然后根据sa实现selpg的功能。  
```go
func processInput(sa selpgArgs) {
	var cmd *exec.Cmd
	var fin io.ReadCloser
	var fout io.WriteCloser
	var lineCtr, // line counter
		pageCtr int // page counter
	var err error
	var line []byte
	if sa.inFilename == "" {
		fin = os.Stdin
	} else {

		if fin, err = os.Open(sa.inFilename); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if sa.printDest == "" {
		fout = os.Stdout
	} else {

		cmd = exec.Command("lp", "-d", sa.printDest)
		cmd.Stderr = os.Stderr
		fout, _ = cmd.StdinPipe()
	}

	finReader := bufio.NewReader(fin)
	//lines-delimited page type
	if sa.pageType == false {
		lineCtr = 0
		pageCtr = 1
		for {
			if line, err = finReader.ReadBytes('\n'); err == io.EOF {
				break
			}
			lineCtr++
			if lineCtr > sa.pageLen {
				pageCtr++
				lineCtr = 1
			}
			if pageCtr >= sa.startPage && pageCtr <= sa.endPage {
				fmt.Fprintf(fout, string(line))
			}
		}
	} else {
		pageCtr = 1
		var c byte
		for {

			if c, err = finReader.ReadByte(); err == io.EOF {
				break
			} else if c == '\f' {
				pageCtr++
			}

			if pageCtr >= sa.startPage && pageCtr <= sa.endPage {
				fmt.Fprintf(fout, string(c))
			}
		}
	}
	if sa.printDest != "" {
		cmd.Run()
	}
	if pageCtr < sa.startPage {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d), no output written\n", progname, sa.startPage, pageCtr)
	} else if pageCtr < sa.endPage {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, sa.endPage, pageCtr)
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}
```


