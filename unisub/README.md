# unisub
Find Unicode characters that might be converted to a given character under some encoding conversions.

## Install
```
▶ go get -u github.com/tomnomnom/hacks/unisub
```

## Usage

```
▶ unisub '@'
＠ U+FF20 %EF%BC%A0
﹫ U+FE6B %EF%B9%AB
```

```
▶ unisub '<'
＜ U+FF1C %EF%BC%9C
﹤ U+FE64 %EF%B9%A4
```

```
▶ unisub 's'
ｓ U+FF53 %EF%BD%93
𝐬 U+1D42C %F0%9D%90%AC
𝑠 U+1D460 %F0%9D%91%A0
𝒔 U+1D494 %F0%9D%92%94
𝓈 U+1D4C8 %F0%9D%93%88
𝓼 U+1D4FC %F0%9D%93%BC
𝔰 U+1D530 %F0%9D%94%B0
𝕤 U+1D564 %F0%9D%95%A4
𝖘 U+1D598 %F0%9D%96%98
𝗌 U+1D5CC %F0%9D%97%8C
𝘀 U+1D600 %F0%9D%98%80
𝘴 U+1D634 %F0%9D%98%B4
𝙨 U+1D668 %F0%9D%99%A8
𝚜 U+1D69C %F0%9D%9A%9C
ⓢ U+24E2 %E2%93%A2
ˢ U+02E2 %CB%A2
ₛ U+209B %E2%82%9B
ſ U+017F %C5%BF
```

