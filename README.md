<img src="https://github.com/facefunk/starlet/blob/master/logo.svg" width="256" height="256" alt="starlet" title="startlet">

# starlet

[![Reference][godoc-img]][godoc-url]
[![Report][report-img]][report-url]
[![License][license-img]][license-url]

Generates CSS from `.strlt` files. Very similar to Stylus, but with higher compression.

## Basic usage

```starlet
body
	color black
	font-size 100%
	padding 1rem
```

## State

```starlet
a
	color blue

	:hover
		color red
```

## Classes

```starlet
a
	color blue

	// "active" elements inside a link
	.active
		color red

	// links that have the "active" class
	&.active
		color red
```

## Multiple selectors

```starlet
// All in one line
h1, h2, h3
	color orange

// Split over multiple lines
h4,
h5,
h6
	color purple
```

## Variables

```starlet
text-color = black
transition-speed = 200ms

body
	font-size 100%
	color text-color

a
	color blue
	transition color transition-speed ease
	
	:hover
		color red
```

## Mixins

```starlet
mixin horizontal
	display flex
	flex-direction row

mixin vertical
	display flex
	flex-direction column
```

Mixins can be used like this:

```starlet
#sidebar
	vertical
```

## Animations

```starlet
animation rotate
	0%
		transform rotateZ(0)
	100%
		transform rotateZ(360deg)

animation pulse
	0%, 100%
		transform scale3D(0.4, 0.4, 0.4)
	50%
		transform scale3D(0.9, 0.9, 0.9)
```

## Quick media queries

```starlet
body
	vertical

> 800px
	body
		horizontal
```

[godoc-img]: https://godoc.org/github.com/facefunk/starlet?status.svg
[godoc-url]: https://godoc.org/github.com/facefunk/starlet
[report-img]: https://goreportcard.com/badge/github.com/facefunk/starlet
[report-url]: https://goreportcard.com/report/github.com/facefunk/starlet
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://github.com/facefunk/starlet/blob/master/LICENSE
