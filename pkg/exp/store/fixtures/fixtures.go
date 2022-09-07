package fixtures

import (
	"io/fs"

	"github.com/nix-community/go-nix/pkg/exp/store/model"
	"github.com/nix-community/go-nix/pkg/exp/store/treestore"
)

//nolint:gochecknoglobals
var (
	BlobEmptyStruct = &model.Blob{
		Contents: []byte{},
	}
	BlobEmptySerialized = []byte{
		0x62, 0x6c, 0x6f, 0x62, // "blob"
		0x20, // space
		0x30, // "0"
		0x00, // null byte
		// empty, no data
	}
	BlobEmptySha1Digest = []byte{
		0xe6, 0x9d, 0xe2, 0x9b, 0xb2, 0xd1, 0xd6, 0x43, 0x4b, 0x8b,
		0x29, 0xae, 0x77, 0x5a, 0xd8, 0xc2, 0xe4, 0x8c, 0x53, 0x91,
	}

	BlobBarStruct = &model.Blob{
		Contents: []byte("Hello World\n"),
	}
	BlobBarSerialized = []byte{
		0x62, 0x6c, 0x6f, 0x62, // "blob"
		0x20,       // space
		0x31, 0x32, // "12"
		0x00,                         // null byte
		0x48, 0x65, 0x6c, 0x6c, 0x6f, // "Hello"
		0x20,                               // space
		0x57, 0x6f, 0x72, 0x6c, 0x64, 0x0a, // "World\n"
	}
	BlobBarSha1Digest = []byte{
		0x55, 0x7d, 0xb0, 0x3d, 0xe9, 0x97, 0xc8, 0x6a, 0x4a, 0x02,
		0x8e, 0x1e, 0xbd, 0x3a, 0x1c, 0xeb, 0x22, 0x5b, 0xe2, 0x38,
	}

	BlobBazStruct = &model.Blob{
		Contents: []byte("foo"),
	}
	BlobBazSerialized = []byte{
		0x62, 0x6c, 0x6f, 0x62, // "blob"
		0x20,             // space
		0x33,             // "3"
		0x00,             // bull byte
		0x66, 0x6f, 0x6f, // "foo"
	}
	BlobBazSha1Digest = []byte{
		0x19, 0x10, 0x28, 0x15, 0x66, 0x3d, 0x23, 0xf8, 0xb7, 0x5a,
		0x47, 0xe7, 0xa0, 0x19, 0x65, 0xdc, 0xdc, 0x96, 0x46, 0x8c,
	}

	BlobFooStruct = &model.Blob{
		Contents: []byte("bar"),
	}
	BlobFooSerialized = []byte{
		0x62, 0x6c, 0x6f, 0x62, // "blob"
		0x20,             // space
		0x33,             // "3"
		0x00,             // null byte
		0x62, 0x61, 0x72, // "bar"
	}
	BlobFooSha1Digest = []byte{
		0xba, 0x0e, 0x16, 0x2e, 0x1c, 0x47, 0x46, 0x9e, 0x3f, 0xe4,
		0xb3, 0x93, 0xa8, 0xbf, 0x8c, 0x56, 0x9f, 0x30, 0x21, 0x16,
	}

	// tree1:
	//
	//  git cat-file -p 79adc4923d0e3d1c620943f58c118368798329d7
	//  040000 tree 29a422c19251aeaeb907175e9b3219a9bed6c616	bab
	//  100644 blob 557db03de997c86a4a028e1ebd3a1ceb225be238	bar
	//  120000 blob 19102815663d23f8b75a47e7a01965dcdc96468c	baz
	//  100755 blob ba0e162e1c47469e3fe4b393a8bf8c569f302116	foo

	Tree1Struct = &model.Tree{
		Entries: []*model.Entry{
			{
				Id: []byte{
					0x29, 0xa4, 0x22, 0xc1, 0x92, 0x51, 0xae, 0xae, 0xb9, 0x07,
					0x17, 0x5e, 0x9b, 0x32, 0x19, 0xa9, 0xbe, 0xd6, 0xc6, 0x16,
				},
				Mode: model.Entry_MODE_DIRECTORY,
				Name: "bab",
			}, {
				Id: []byte{
					0x55, 0x7d, 0xb0, 0x3d, 0xe9, 0x97, 0xc8, 0x6a, 0x4a, 0x02,
					0x8e, 0x1e, 0xbd, 0x3a, 0x1c, 0xeb, 0x22, 0x5b, 0xe2, 0x38,
				},
				Mode: model.Entry_MODE_FILE_REGULAR,
				Name: "bar",
			}, {
				Id: []byte{
					0x19, 0x10, 0x28, 0x15, 0x66, 0x3d, 0x23, 0xf8, 0xb7, 0x5a,
					0x47, 0xe7, 0xa0, 0x19, 0x65, 0xdc, 0xdc, 0x96, 0x46, 0x8c,
				},
				Mode: model.Entry_MODE_SYMLINK,
				Name: "baz",
			}, {
				Id: []byte{
					0xba, 0x0e, 0x16, 0x2e, 0x1c, 0x47, 0x46, 0x9e, 0x3f, 0xe4,
					0xb3, 0x93, 0xa8, 0xbf, 0x8c, 0x56, 0x9f, 0x30, 0x21, 0x16,
				},
				Mode: model.Entry_MODE_FILE_EXECUTABLE,
				Name: "foo",
			},
		},
	}

	Tree1Serialized = []byte{
		0x74, 0x72, 0x65, 0x65, // "tree"
		0x20,             // space
		0x31, 0x32, 0x33, // content size (123)
		0x00, // null byte

		// first entry
		0x34, 0x30, 0x30, 0x30, 0x30, // 40000 (type directory)
		0x20,             // space
		0x62, 0x61, 0x62, // "bab" (name)
		0x00, // null byte
		0x29, 0xa4, 0x22, 0xc1, 0x92, 0x51, 0xae, 0xae, 0xb9, 0x07,
		0x17, 0x5e, 0x9b, 0x32, 0x19, 0xa9, 0xbe, 0xd6, 0xc6, 0x16, // hash (20 bytes for sha1)

		// second entry
		0x31, 0x30, 0x30, 0x36, 0x34, 0x34, // 100644 (type regular)
		0x20,             // space
		0x62, 0x61, 0x72, // "bar" (name)
		0x00, // null byte
		0x55, 0x7d, 0xb0, 0x3d, 0xe9, 0x97, 0xc8, 0x6a, 0x4a, 0x02,
		0x8e, 0x1e, 0xbd, 0x3a, 0x1c, 0xeb, 0x22, 0x5b, 0xe2, 0x38, // hash (20 bytes for sha1)

		// third entry
		0x31, 0x32, 0x30, 0x30, 0x30, 0x30, // 120000 (type symlink)
		0x20,             // space
		0x62, 0x61, 0x7a, // "baz" (name)
		0x00, // null byte
		0x19, 0x10, 0x28, 0x15, 0x66, 0x3d, 0x23, 0xf8, 0xb7, 0x5a,
		0x47, 0xe7, 0xa0, 0x19, 0x65, 0xdc, 0xdc, 0x96, 0x46, 0x8c, // hash (20 bytes for sha1)

		// fourth entry
		0x31, 0x30, 0x30, 0x37, 0x35, 0x35, // 100755 (type executable)
		0x20,             // space
		0x66, 0x6f, 0x6f, // "foo" (name)
		0x00, // null byte
		0xba, 0x0e, 0x16, 0x2e, 0x1c, 0x47, 0x46, 0x9e, 0x3f, 0xe4,
		0xb3, 0x93, 0xa8, 0xbf, 0x8c, 0x56, 0x9f, 0x30, 0x21, 0x16, // hash (20 bytes for sha1)
	}

	Tree1Sha1Digest = []byte{
		0x79, 0xad, 0xc4, 0x92, 0x3d, 0xe, 0x3d, 0x1c, 0x62, 0x9,
		0x43, 0xf5, 0x8c, 0x11, 0x83, 0x68, 0x79, 0x83, 0x29, 0xd7,
	}

	// tree2:
	//  git cat-file -p 29a422c19251aeaeb907175e9b3219a9bed6c616
	//  100644 blob e69de29bb2d1d6434b8b29ae775ad8c2e48c5391	.keep

	Tree2Struct = &model.Tree{
		Entries: []*model.Entry{
			{
				Id: []byte{
					0xe6, 0x9d, 0xe2, 0x9b, 0xb2, 0xd1, 0xd6, 0x43, 0x4b, 0x8b,
					0x29, 0xae, 0x77, 0x5a, 0xd8, 0xc2, 0xe4, 0x8c, 0x53, 0x91,
				},
				Mode: model.Entry_MODE_FILE_REGULAR,
				Name: ".keep",
			},
		},
	}

	Tree2Serialized = []byte{
		0x74, 0x72, 0x65, 0x65, // "tree"
		0x20,       // space
		0x33, 0x33, // content size (33)
		0x00, // null byte

		// first entry
		0x31, 0x30, 0x30, 0x36, 0x34, 0x34, // 100644 (type regular)
		0x20, // space
		0x2e, 0x6b, 0x65, 0x65, 0x70,
		0x00, // null byte
		0xe6, 0x9d, 0xe2, 0x9b, 0xb2, 0xd1, 0xd6, 0x43, 0x4b, 0x8b,
		0x29, 0xae, 0x77, 0x5a, 0xd8, 0xc2, 0xe4, 0x8c, 0x53, 0x91, // hash (20 bytes for sha1)
	}

	Tree2Sha1Digest = []byte{
		0x29, 0xa4, 0x22, 0xc1, 0x92, 0x51, 0xae, 0xae, 0xb9, 0x7,
		0x17, 0x5e, 0x9b, 0x32, 0x19, 0xa9, 0xbe, 0xd6, 0xc6, 0x16,
	}

	Tree2Entries = []treestore.Entry{
		{
			Path:     "/",
			DirEntry: NewMockDirEntry("/", 0, fs.ModePerm|fs.ModeDir),
		}, {
			ID:       BlobEmptySha1Digest,
			Path:     "/.keep",
			DirEntry: NewMockDirEntry(".keep", 0, 0o644),
		},
	}

	WholeTreeEntries = []treestore.Entry{
		{
			Path:     "/",
			DirEntry: NewMockDirEntry("/", 0, fs.ModePerm|fs.ModeDir),
		}, {
			ID:       Tree2Sha1Digest,
			Path:     "/bab",
			DirEntry: NewMockDirEntry("bab", 0, fs.ModePerm|fs.ModeDir),
		}, {
			ID:       BlobEmptySha1Digest,
			Path:     "/bab/.keep",
			DirEntry: NewMockDirEntry(".keep", 0, 0o644),
		}, {
			ID:       BlobBarSha1Digest,
			Path:     "/bar",
			DirEntry: NewMockDirEntry("bar", 12, 0o644),
		}, {
			ID:       BlobBazSha1Digest,
			Path:     "/baz",
			DirEntry: NewMockDirEntry("baz", 3, fs.ModePerm|fs.ModeSymlink),
		}, {
			ID:       BlobFooSha1Digest,
			Path:     "/foo",
			DirEntry: NewMockDirEntry("foo", 3, 0o700),
		},
	}
)
