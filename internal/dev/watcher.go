// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package dev

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

// Watch the directory in the path, and send a signal when it changes.
type watcher struct {
	path    string
	changed chan<- struct{}
}

func (w *watcher) start(ctx context.Context, ready chan<- struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(w.path)
	if err != nil {
		return err
	}

	logger.Print("fs: watching changes in ", w.path)
	ready <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) {
				select {
				case w.changed <- struct{}{}:
					logger.Print("fs: directory changed")
				default:
				}
			}
		case err := <-watcher.Errors:
			logger.Print("fs: ", err)
		}
	}
}
