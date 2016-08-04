#
# ASK STATUS
#

$iisInfo = Get-ItemProperty HKLM:\SOFTWARE\Microsoft\InetStp\
$version = [decimal]"$($iisInfo.MajorVersion).$($iisInfo.MinorVersion)"

#
# output
#

Write-Host $version
exit 0

