/*
Package replay is a simple osu! replay parser made in Go. Still undergoing development, expect breaking changes.

Information about .osr files can be found on the official wiki.
https://osu.ppy.sh/wiki/en/osu%21_File_Formats/Osr_%28file_format%29

To decode a replay simply pass its file representation as an io.Reader into replay.New:

	r := replay.New(readerOfOsrFile)
	fmt.Printf("%s played this replay", r.Name)

You now have a replay struct that you can analyse and edit. Have fun.
CLI tool coming soon.
 */
package replay
