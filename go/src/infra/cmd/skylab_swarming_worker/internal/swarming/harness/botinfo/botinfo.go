// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package botinfo implements opening and closing a bot's botinfo stored on
// disk.
package botinfo

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"go.chromium.org/luci/common/errors"

	"infra/cmd/skylab_swarming_worker/internal/botinfo"
	"infra/cmd/skylab_swarming_worker/internal/swarming"
)

// Store holds a bot's botinfo and adds a Close method.
type Store struct {
	botinfo.BotInfo
	bot *swarming.Bot
}

// Close writes the BotInfo back to disk.  This method does nothing on
// subsequent calls.  This method is safe to call on a nil pointer.
func (s *Store) Close() error {
	if s == nil {
		return nil
	}
	if s.bot == nil {
		return nil
	}
	data, err := botinfo.Marshal(&s.BotInfo)
	if err != nil {
		return errors.Annotate(err, "close botinfo").Err()
	}
	if err := ioutil.WriteFile(botinfoFilePath(s.bot), data, 0666); err != nil {
		return errors.Annotate(err, "close botinfo").Err()
	}
	s.bot = nil
	return nil
}

// Open loads the BotInfo for the Bot.  The BotInfo should be closed
// afterward to write it back.
func Open(b *swarming.Bot) (*Store, error) {
	s := Store{bot: b}
	data, err := ioutil.ReadFile(botinfoFilePath(b))
	if err != nil {
		return nil, errors.Annotate(err, "open botinfo").Err()
	}
	if err := botinfo.Unmarshal(data, &s.BotInfo); err != nil {
		return nil, errors.Annotate(err, "open botinfo").Err()
	}
	return &s, nil
}

// botinfoFilePath returns the path for caching dimensions for the given bot.
func botinfoFilePath(b *swarming.Bot) string {
	return filepath.Join(botinfoDirPath(b), fmt.Sprintf("%s.json", b.DUTID))
}

// botinfoDir returns the path to the cache directory for the given bot.
func botinfoDirPath(b *swarming.Bot) string {
	return filepath.Join(b.AutotestPath, "swarming_state")
}