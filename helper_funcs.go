package main

import "muhammaddev/internal/database"

func dbPlaylistToPlaylist(dbPlaylist database.Playlist) Playlist {
	return Playlist{
		ID:          dbPlaylist.ID.String(),
		Name:        dbPlaylist.Name,
		Description: dbPlaylist.Description.String,
	}
}

func dbPlaylistsToPlaylists(dbPlaylists []database.Playlist) []Playlist {
	playlists := []Playlist{}
	for _, dbPlaylist := range dbPlaylists {
		playlists = append(playlists, dbPlaylistToPlaylist(dbPlaylist))
	}

	return playlists
}

func dbTutorialToTutorial(dbTutorial database.Tutorial) Tutorial {
	return Tutorial{
		ID:          dbTutorial.ID.String(),
		Title:       dbTutorial.Title,
		TutorialUrl: dbTutorial.TutorialUrl,
		Description: dbTutorial.Description,
		YoutubeLink: dbTutorial.YoutubeLink,
		PlaylistID:  dbTutorial.PlaylistID.String(),
	}
}

func dbTutorialsToTutorials(dbTutorials []database.Tutorial) []Tutorial {
	tutorials := []Tutorial{}
	for _, dbTutorial := range dbTutorials {
		tutorials = append(tutorials, dbTutorialToTutorial(dbTutorial))
	}

	return tutorials
}

func dbImageToImage(dbImage database.Image) Image {
	return Image{
		ID:       dbImage.ID.String(),
		ImageUrl: dbImage.ImageUrl,
	}
}

func dbImagesToImages(dbImages []database.Image) []Image {
	images := []Image{}
	for _, image := range dbImages {
		images = append(images, dbImageToImage(image))
	}
	return images
}
