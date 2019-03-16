#!/usr/bin/perl
# Assumes a list of files with the format (from "jhead"):
#
#/data/mp3/Music/Mine/Adema - Adema/folder.jpg
#File name    : /data/mp3/Music/Mine/Adema - Adema/folder.jpg
#File size    : 25559 bytes
#File date    : 2005:01:01 15:12:54
#Resolution   : 300 x 298
#
#
# Created by:
# (OLD) find /data/mp3/Music/Net -name folder.jpg -print -exec jhead {} ";"
# find /data/media/mp3 -name folder.jpg -print -exec jhead {} ";" > ToGain.txt

$Filename="ToGain.txt";
$Tolerance=.2; # .8 - 1.2
$Sides=500; # Size of X and Y dimensions
$PBM="/usr/bin/";
$Jhead="/usr/bin/jhead";

# Populate @myfile with the contents of the file
open FOO,$Filename or die $!;
while(<FOO>)
{
	unshift @myfile, $_;
}
close FOO;

# Build list of filenames
$i=0;
for (@myfile)
{
	if (/File name/o)
	{
		$mytempname=substr($_,15); # Strip off "File name :"
		chomp($mytempname);
		$myname[$i]=quotemeta($mytempname); # Escape characters
		$i++;
	}
}

# And build list of resolutions
$i=0;
for (@myfile)
{
	if (/Resolution/o)
	{
		@mysres=split(' ',$_); # This is the "Resolution" string
		@myxres[$i]=$mysres[2]; # X is the 2nd argument
		@myyres[$i]=$mysres[4]; # Y is the 4th argument
		$i++;
	}
}

$i=$#myname;
while ($i >= 0)
{
	#print "X=$myxres[$i] Y=$myyres[$i]\n";
	$ratio=$myxres[$i]/$myyres[$i];
	#print "Ratio=$ratio\n";

	# If the picture is "about square"
	# and the picture is not already $Sides x $Sides
	# then resize it.
	if ( (($ratio < 1+$Tolerance) && ($ratio > 1-$Tolerance)) && (($myxres[$i] ne $Sides) || ($myyres[$i] ne $Sides)) )
	{
		print "Name=$myname[$i]\n";
		#print "Ratio=$ratio\n";
		`$PBM/jpegtopnm $myname[$i] | $PBM/pnmscale -xsize=500 -ysize=500 | $PBM/pnmtojpeg -quality=60 -optimize >/tmp/folder_tmp.jpg`;
		`mv -f /tmp/folder_tmp.jpg $myname[$i]`;
	}

	`$Jhead -purejpg $myname[$i]`; # Strip all non-picture elements

	$i--;
}

