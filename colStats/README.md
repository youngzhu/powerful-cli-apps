## v1
通过profile工具发现，分配的内存较大且频繁。是`csv.ReadAll()`的缘故。
其实不用一次性把文件内容全部读入内存，可以一条条读，读一条处理一条。
优化后，CPU时间利用率提高了40%左右。

## v2
通过trace工具发现，只使用了一个CPU，没有利用多核。
每个文件产生一个协程，上万个协程，反而更慢了。

## v3
控制协程的数量。