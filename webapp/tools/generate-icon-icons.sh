#!/bin/bash

INK="/usr/bin/inkscape"

if [ "$(uname)" == "Darwin" ]; then
    INK="/Applications/Inkscape.app/Contents/MacOS/inkscape"
fi

if [[ -z "$1" ]] 
then
	echo "SVG file needed."
	exit;
fi

SVG="$1"

OUT="ios/App/App/Assets.xcassets/AppIcon.appiconset/AppIcon"

$INK -D -o "$OUT-20x20@1x.png" -w 20 -h 20 $SVG
$INK -D -o "$OUT-20x20@2x-1.png" -w 40 -h 40 $SVG
$INK -D -o "$OUT-20x20@2x.png" -w 40 -h 40 $SVG
$INK -D -o "$OUT-20x20@3x.png" -w 60 -h 60 $SVG
$INK -D -o "$OUT-29x29@1x.png" -w 29 -h 29 $SVG
$INK -D -o "$OUT-29x29@2x-1.png" -w 58 -h 58 $SVG
$INK -D -o "$OUT-29x29@2x.png" -w 58 -h 58 $SVG
$INK -D -o "$OUT-29x29@3x.png" -w 87 -h 87 $SVG
$INK -D -o "$OUT-40x40@1x.png" -w 40 -h 40 $SVG
$INK -D -o "$OUT-40x40@2x-1.png" -w 80 -h 80 $SVG
$INK -D -o "$OUT-40x40@2x.png" -w 80 -h 80 $SVG
$INK -D -o "$OUT-40x40@3x.png" -w 120 -h 120 $SVG
$INK -D -o "$OUT-60x60@2x.png" -w 120 -h 120 $SVG
$INK -D -o "$OUT-60x60@3x.png" -w 180 -h 180 $SVG
$INK -D -o "$OUT-76x76@1x.png" -w 76 -h 76 $SVG
$INK -D -o "$OUT-76x76@2x.png" -w 152 -h 152 $SVG
$INK -D -o "$OUT-83.5x83.5@2x.png" -w 167 -h 167 $SVG
$INK -D -o "$OUT-512@2x.png" -w 1024 -h 1024 $SVG
