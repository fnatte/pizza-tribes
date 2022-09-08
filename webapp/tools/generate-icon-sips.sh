#!/bin/bash

SIPS="/usr/bin/sips"

if [[ -z "$1" ]] 
then
	echo "PNG file needed."
	exit;
fi

PNG="$1"

OUT="ios/App/App/Assets.xcassets/AppIcon.appiconset/AppIcon"

$SIPS -o "$OUT-20x20@1x.png" -Z 20 $PNG
$SIPS -o "$OUT-20x20@2x-1.png" -Z 40 $PNG
$SIPS -o "$OUT-20x20@2x.png" -Z 40 $PNG
$SIPS -o "$OUT-20x20@3x.png" -Z 60 $PNG
$SIPS -o "$OUT-29x29@1x.png" -Z 29 $PNG
$SIPS -o "$OUT-29x29@2x-1.png" -Z 58 $PNG
$SIPS -o "$OUT-29x29@2x.png" -Z 58 $PNG
$SIPS -o "$OUT-29x29@3x.png" -Z 87 $PNG
$SIPS -o "$OUT-40x40@1x.png" -Z 40 $PNG
$SIPS -o "$OUT-40x40@2x-1.png" -Z 80 $PNG
$SIPS -o "$OUT-40x40@2x.png" -Z 80 $PNG
$SIPS -o "$OUT-40x40@3x.png" -Z 120 $PNG
$SIPS -o "$OUT-60x60@2x.png" -Z 120 $PNG
$SIPS -o "$OUT-60x60@3x.png" -Z 180 $PNG
$SIPS -o "$OUT-76x76@1x.png" -Z 76 $PNG
$SIPS -o "$OUT-76x76@2x.png" -Z 152 $PNG
$SIPS -o "$OUT-83.5x83.5@2x.png" -Z 167 $PNG
$SIPS -o "$OUT-512@2x.png" -Z 1024 $PNG
