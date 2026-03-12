package db

import (
    "time"

    "github.com/google/uuid"
)

// NOTE: For people wondering about use of `panic` in String() methods, read the middleware package, it already
// handles server panics gracefully by restarting the server and sending a 500 error to the client. While this
// may not be the best way to handle this sort of errors in a critical high volume server but this is supposed to
// be a local server and so panicing will leave logs full of red color (with the help of the custom logger) and
// will be easier to debug and find issues early. Just make sure to do all the DB ops in transaction to avoid any
// potential data loss if the server was processing something with the database so that rollbacks are possible.

// INFO: Two type structures exists for Sakuhin, one represents data before or at DB layer and one respresents
// data after leaving DB layer which is to be used by the entire application. Also note that data
// validation will take place in database with sql schema designed to prevent incorrect data from
// even entering the database in the first place so that we don't need an app layer validator.

// Schema before or at DB layer: `db.Sakuhin.PrequelID -> uuid.UUID`. (foreign key to table for prequel) [YOU ARE HERE]
// Schema after DB layer, to be used by the app: `types.Sakuhin.Prequel - types.Sakuhin` (the entire table for prequel)
// Depth and if we should omit inclusion are controllable while data fetching from DB.

// Sakuhin is designed with Manga and Doujin in mind, it can probably loosely represent other work such as LN too.
type Sakuhin struct {
    // ID is primary key of a row that is used to uniquely identify it.
    ID uuid.UUID // Primary key
    // Title is the title of work the row represents in two supported language, note that I do not plan on adding
    // support for any new language as I made this service with my own use case in mind and I'm proficient in English
    // and am currently learning Japanese.
    Title StringWithLang
    // Description is the description of the work, I don't think this needs much explanation.
    Description StringWithLang
    // CoverArts will store the front page of the work in case it is provided, if not provided we use the first page
    // of the work. I will use an unique ID in a very specific directory where all the cover arts will be stored, this
    // way we can change the path to directory any time we want and use the ID to get the file within that specific
    // directory.
    CoverArts *[]uuid.UUID
    // Tag is probably gonna be a new table with a foreign key to the row it represents, it contains a bunch of tags.
    TagsID uuid.UUID // Foreign key to `SakuhinTag`
    // I'm not sure if category will come in any use but let's just keep it just in case to distinguish.
    Category SakuhinCategory
    // Demographic is the target audience base of a specific work. Can be nil in some cases.
    Demographic   *SakuhinDemographic  // can be nil if unspecified demographic or multi-demographic
    ContentRating SakuhinContentRating // self explanatory
    // FIXME: Study junction table and implement one
    Artist []uuid.UUID // Foreign key (kind of, foreign key cannot be an array) to `SakuhinCreator`
    Author []uuid.UUID // Foreign key (kind of, foreign key cannot be an array) to `SakuhinCreator`
    // All the recognizable character in the work, this is probably also a new table with a foreign key to the row it
    // represents.
    // FIXME: Study junction table and implement one
    Character []uuid.UUID // Foreign key (kind of, foreign key cannot be an array) to `SakuhinCharacter`
    // Language of the work, for most cases it's gonna be en or jp but there are some doujin that are only available
    // online in language like cn or kr and is often left without translation so I made language a generic string,
    // it would be better if I can change this to a enum later in the future, note that the value must be the
    // language name in full, example: "english", "japanese", etc.
    Language string
    // ID of a specific series, this is unique only between different series and a series with multiple volumes, cour,
    // parts, etc. regardless of how it is being released/serialized, would have the same SeriesID.
    SeriesID uuid.UUID // All installments/parts of a series share this same ID
    // IDs of the row that is prequel to this installment. Note that while the field/property says prequel, it is not
    // prequel in the sense of storyline/timeline but rather the release date instead.
    // FIXME: Study junction table and implement one
    PrequelID *[]uuid.UUID
    // Same as previously defined field PrequelID just for Sequels and as mentioned previously, it is also based on
    // release date rather than the timeline of the series as there are certain series such as fate where it is hard to
    // say which specific installment is prequel and which installment is sequel due to it sheer amount of references to
    // different series and the sheer number of installments that it currently has.
    // FIXME: Study junction table and implement one
    SequelID *[]uuid.UUID
    // Parodies is almost certainly dedicated to doujin and is pretty self explanatory.
    // FIXME: Study junction table and implement one
    ParodiesID *[]uuid.UUID // Foreign key (kind of, foreign key cannot be an array) to `SakuhinParodies` // if nil then the series is an original series
    PageCount  uint16       // self explanatory
    // Seriel number of the installment release count starting from 1..n installments.
    InstallmentNumber uint32 // Is there any series with more than 4 billion installments?
    // This field was mostly for the frontend but I guess we can use it to see if the installment is a "special", "extra", etc.
    DisplayLabel DisplayLabel
    ReleasedAt   time.Time // self explanatory
    UploadedAt   time.Time // self explanatory
    UpdatedAt    time.Time // self explanatory
}

// Modern language don't have unions so I'm just gonna enforce in the frontend that a title must be provided, either En or Jp.
type StringWithLang struct {
    English  *string `json:"english"`
    Japanese *string `json:"japanese"`
}

type SakuhinTag struct {
    ID uuid.UUID // Primary key
    // TODO: Probably better to have a enum dedicated to themes and genres defining all possible combinations
    Themes  []string  `json:"themes"`  // Themes are stuff like "foot fetish", "tomboy", "gore", "ugly bastard", "slasher", etc.
    Genres  []string  `json:"genres"`  // Genres are stuff like "Horror", "Romance", "Yuri", "Isekai", "Comedy", "Action", etc.
    Circles *[]string `json:"circles"` // Circle is most likely gonna be nil for manga and list of circles/groups for doujin.
}

// Role agnostic strucutre which can be used to represent either an artist, author, etc.
// This is probably gonna be a table which is autogenerated whena new author/artist name is referenced while uploading a work.
// This way we can add autocomplete feature displaying the name when the user types while uploading new works in the frontend UI.
type SakuhinCreator struct {
    ID   uuid.UUID      `json:"id"` // Primary key
    Name StringWithLang `json:"name"`
}

// Separate table
type SakuhinRepresentation struct {
    ID           uuid.UUID `json:"id"`            // Primary key
    SeriesID     uuid.UUID `json:"series_id"`     // Foreign key to `Sakuhin.SeriesID`
    ChapterCount *uint64   `json:"chapter_count"` // will be nil in cases where we cannot count chapter such as in tankouban release, etc.
    // FIXME: Improve the volumes by making it a table or decide if we want to keep it as stringified JSON.
    Volumes *[]struct {
        SerialNumber uint64 `json:"serial_number"`
        // one way to index volume is with start and end installment number
        StartInstallment   *uint64   `json:"start_installment"`
        ExcludeInstallment *[]uint64 `json:"exclude_installment"` // if an installment released between two points is a special, spin off, etc. we should exclude
        EndInstallment     *uint64   `json:"end_installment"`
        // another way to index volume is with an array of series id, can be usefull if the volume is a compilation of doujinshis.
        // the way most online website index a compilation is probably using a single cbz file for the entire upload of the volume.
        // we can instead upload the parts/installments/works separetely and then index them together as a single volume/compilation
        // this is very convienient as we can either read a single work or a compilation of works of a circle (assuming that a volume contains works from same circle).
        SeriesIDs *[]uuid.UUID `json:"series_ids"`
    } `json:"volumes"` // will be nil in cases where the work is not supposed to be counted as a "volume", example a single doujinshi vs a compilation of doujinshis.
}

type SakuhinCharacter struct {
    ID            uuid.UUID    `json:"id"` // Primary key
    Role          string       `json:"role"`
    EnglishTitle  *string      `json:"english_title"`  // if nil, then character name is probably not translated
    JapaneseTitle *string      `json:"japanese_title"` // if nil, then character is probably not of japanese origin
    Photos        *[]uuid.UUID `json:"photos"`         // if nil, then the character probably is a mythical faceless character or just not been revealed in the story yet
}

type SakuhinParodies struct {
    ID         uuid.UUID      `json:"id"` // primary key
    Title      StringWithLang `json:"title"`
    Characters []uuid.UUID    `json:"characters"`
}

// Separate table
type SakuhinPublicationStat struct {
    SeriesID        uuid.UUID `json:"series_id"` // foreign key
    InitialRelease  time.Time `json:"initial_release"`
    LatestRelease   time.Time `json:"latest_release"`
    CurrentStatus   Status    `json:"status"`
    ReleaseSchedule Schedule  `json:"schedule"`
}

