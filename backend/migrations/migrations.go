package migrations

import "embed"

// Files embeds all versioned SQL migration files in lexicographical order.
//
//go:embed *.sql
var Files embed.FS
