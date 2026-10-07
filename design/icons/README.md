# Totipo icon system

This directory is the canonical source for Totipo application icon artwork and derived platform assets.

## Canonical mark

The Totipo mark is deliberately simple:

- one circular ring;
- exactly **three timer ticks**;
- exactly **three token dots**;
- **no lock or shackle** at any size.

All platforms use the same mark. Do not add a lock, sync arrow, keyhole, extra ticks, extra dots, text, or other feature-specific imagery to platform variants.

## Source-of-truth files

- `totipo-app-icon.svg` — canonical precomposed application icon: white Totipo mark on the approved green gradient tile.
- `totipo-mark.svg` — canonical symbol geometry with `currentColor`, for contexts that supply their own foreground/background treatment.
- `totipo-mark-white.svg` — the same mark with a fixed white foreground, convenient for launcher/adaptive-icon generation.

The SVG files above are authoritative. Files below `reference/` and `platforms/` are derived exports.

## Icon color

The app-icon artwork uses the approved green gradient from the master SVG:

- `#46FB70`
- `#08D267`
- `#028B55`
- `#026344`

These are **brand artwork colors**, not the application's semantic UI palette. Warning, error, destructive, and other semantic states should continue to use their design-system colors rather than repurposing the icon green.

## Size policy

The simplified mark is used unchanged at all sizes. There is no separate micro logo and no feature removal at small sizes.

Reference PNGs are committed at 16, 24, 32, 48, 64, 128, 256, 512, and 1024 pixels under `reference/png/`.

When a platform can consume vector or adaptive artwork, prefer the canonical SVG geometry over scaling a small raster.

## Platform assets

- `platforms/android/` — adaptive-icon foreground/background resources, density-specific legacy launcher icons, and Play Store artwork.
- `platforms/ios/AppIcon.appiconset/` — Xcode app-icon asset catalog with `Contents.json`.
- `platforms/macos/` — `.iconset` sources and `.icns` export.
- `platforms/windows/totipo.ico` — multi-resolution Windows icon.
- `platforms/linux/hicolor/` — standard hicolor icon hierarchy plus scalable SVG.
- `platforms/web/` — favicon, Apple touch icon, browser tile, and web manifest assets.

Platform repositories should copy or generate their local assets from this directory rather than maintaining independently drawn versions of the Totipo mark.

## Usage in specification/design documents

For a rendered icon in Markdown, a relative reference can use, for example:

```markdown
![Totipo app icon](icons/reference/png/totipo-256.png)
```

when the document containing the link lives directly under `design/`. Adjust the relative path if the design document is elsewhere.

For diagrams or web-native documentation, prefer `totipo-app-icon.svg` when the renderer supports SVG.

## Change control

A change to any of the following is an icon-design change and should be reviewed explicitly:

- ring geometry;
- tick count or placement;
- dot count or placement;
- icon gradient;
- relative scale/centering of the mark;
- platform-specific deviations from the canonical mark.

Derived exports may be regenerated without changing the icon design as long as they render the canonical master faithfully.
