package main

import "muhammaddev/internal/database"

func dbPlaylistToPlaylist(dbPlaylist database.Playlist) Playlist {
	return Playlist{
		ID:          dbPlaylist.ID,
		Name:        dbPlaylist.Name,
		Description: dbPlaylist.Description,
	}
}
