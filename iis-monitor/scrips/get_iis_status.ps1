# Shell arguments
#
[CmdletBinding()]
Param(
   [Parameter(Mandatory=$True,Position=1)]
   [string]$server,
   [Parameter(Mandatory=$True,Position=2)]
   [string]$website,
   [Parameter(Mandatory=$True,Position=3)]
   [string]$counterlist
   )

Set-Variable OK 0 -option Constant
Set-Variable UNKNOWN 1 -option Constant


#
# ASK STATUS
#

$counter = Get-Counter "\\$server\Web Service($website)\$counterlist"  




# output

$resultstring='UNKNOWN ' + $website + ' or UNKNOWN ' + $counterlist + ' not found' 
$exit_code = $UNKNOWN
  
if ($counter -ne $null) {

	$connections=$counter.CounterSamples.CookedValue

	$resultstring= $connections
	$exit_code = 0
}



Write-Host $resultstring
exit $exit_code

