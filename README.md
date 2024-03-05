# BeatmapPacker

CLI based downloader for osu! standard beatmap packs written in [Go](https://go.dev) (golang).

Outputs into /BeatmapMegapack containing .osz files from the downloaded beatmap packs allowing for easy drag-and-drop to your client.

## How to use

```pwsh
.\beatmappacker.exe <start> <end>
```

Start: Beginning range of pack number to start downloading from.

End: End range of pack number to download.

Example:

```pwsh
.\beatmappacker.exe 1 1348
```

This will download beatmap packs from 1 - 1348.

## Warning

**Very slow!**

osu! file servers do not like downloading multiple files at the same time, highly dependent on your machine's internet speed. The packs are 200MB each, be wary of how much you're downloading.

Prone to breakage as pack links are inconsistent/can change anytime, I'll try to add edge-cases when I can.
