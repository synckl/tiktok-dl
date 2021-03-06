package utils

import (
	models "../models"
	config "../models/config"
	fileio "./fileio"
	log "./log"
)

// IsItemInArchive - Checks if the item is already archived
func IsItemInArchive(upload models.Upload) bool {
	if len(RemoveArchivedItems([]models.Upload{upload})) == 0 {
		return true
	}
	return false
}

// RemoveArchivedItems - Returns items slice without archived items
func RemoveArchivedItems(uploads []models.Upload) []models.Upload {
	archiveFilePath := config.Config.ArchiveFilePath

	if archiveFilePath == "" || !fileio.CheckIfExists(archiveFilePath) {
		return uploads
	}

	removeArchivedItemsDelegate := func(archivedItem string) {
		for i, upload := range uploads {
			if upload.GetUploadID() == archivedItem {
				uploads = append(uploads[:i], uploads[i+1:]...)
			}
		}
	}

	lenBeforeRemoval := len(uploads)
	fileio.ReadFileLineByLine(archiveFilePath, removeArchivedItemsDelegate)

	removedCount := lenBeforeRemoval - len(uploads)
	if removedCount > 0 {
		log.Logf("%d items, found in archive. Skipping...\n", removedCount)
	}

	return uploads
}

// AddItemToArchive - Adds item to archived list
func AddItemToArchive(uploadID string) {
	archiveFilePath := config.Config.ArchiveFilePath

	if archiveFilePath == "" {
		return
	}

	fileio.AppendToFile(uploadID, archiveFilePath)
}
