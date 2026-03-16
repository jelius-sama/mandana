package types

import (
    "mandana/db"
    "time"

    "github.com/google/uuid"
)

// NOTE: For people wondering about use of `panic` in String() methods, read the middleware package, it already
// handles server panics gracefully by restarting the server and sending a 500 error to the client. While this
// may not be the best way to handle this sort of errors in a critical high volume server but this is supposed to
// be a local server and so panicing will leave logs full of red color (with the help of the custom logger) and
// will be easier to debug and find issues early. Just make sure to do all the DB ops in transaction to avoid any
// potential data loss if the server was processing something with the database so that rollbacks are possible.

// TODO: Make two schema type structures, one represents data before or at database level and one respresents
// data after leaving database package which is to be used by the entire application. Also note that data
// validation will take place in database layer with sql statements designed to prevent incorrect data from
// even entering the database in the first place so that we don't need an app layer validator.
// Schema before or at DB layer: `db.Sakuhin.PrequelID - uuid.UUID`.
// Schema after DB layer, to be used by the app: `types.Sakuhin.Prequel - types.Sakuhin`
// Basically foreign key will be replaced by the row that the key points to so that we have the data ready to use.

// TODO: There are some parts that says uuid.UUID but those should be the struct type instead, it's only uuid in
// the database, after it croses the DB package layer it is not gonna be uuid but rather the whole table which
// would be in the struct type of that field, example include `Sakuhin.PrequelID`, `Sakuhin.SequelID`, etc.

// NOTE: We need to define a property that will tell us how deep a recursive field goes or if it is nil.
// Example: `Sakuhin.Prequel` contains `Sakuhin` in an array and if it goes on forever the data will get quite big
// to transfer over the wire therefore we need pagination kind of feature to implement.

// Sakuhin is designed with Manga and Doujin in mind, it can probably loosely represent other work such as LN too.
type Sakuhin struct {
    // ID is primary key of a row that is used to uniquely identify it.
    ID uuid.UUID `json:"id"`
    // Title is the title of work the row represents in two supported language, note that I do not plan on adding
    // support for any new language as I made this service with my own use case in mind and I'm proficient in English
    // and am currently learning Japanese.
    Title StringWithLang `json:"title"`
    // Description is the description of the work, I don't think this needs much explanation.
    Description StringWithLang `json:"description"`
    // CoverArts will store the front page of the work in case it is provided, if not provided we use the first page
    // of the work. I will use an unique ID in a very specific directory where all the cover arts will be stored, this
    // way we can change the path to directory any time we want and use the ID to get the file within that specific
    // directory.
    CoverArts []string `json:"cover_arts"` // compute actual url from uuid in DB layer
    // Tag is probably gonna be a new table with a foreign key to the row it represents, it contains a bunch of tags.
    Tags SakuhinTag `json:"tags"`
    // I'm not sure if category will come in any use but let's just keep it just in case to distinguish.
    Category db.SakuhinCategory `json:"category"`
    // Demographic is the target audience base of a specific work. Can be nil in some cases.
    Demographic   *db.SakuhinDemographic  `json:"demographic"`    // can be nil if unspecified demographic or multi-demographic
    ContentRating db.SakuhinContentRating `json:"content_rating"` // self explanatory
    Artist        []SakuhinCreator        `json:"artist"`         // self explanatory
    Author        []SakuhinCreator        `json:"author"`         // self explanatory
    // All the recognizable character in the work, this is probably also a new table with a foreign key to the row it
    // represents.
    Character SakuhinCharacter `json:"character"`
    // Language of the work, for most cases it's gonna be en or jp but there are some doujin that are only available
    // online in language like cn or kr and is often left without translation so I made language a generic string,
    // it would be better if I can change this to a enum later in the future, note that the value must be the
    // language name in full, example: "english", "japanese", etc.
    Language string `json:"language"`
    // ID of a specific series, this is unique only between different series and a series with multiple volumes, cour,
    // parts, etc. regardless of how it is being released/serialized, would have the same SeriesID.
    SeriesID uuid.UUID `json:"series_id"` // All installments/parts of a series share this same ID
    // IDs of the row that is prequel to this installment. Note that while the field/property says prequel, it is not
    // prequel in the sense of storyline/timeline but rather the release date instead.
    Prequel *[]Sakuhin `json:"prequel_id"`
    // Same as previously defined field PrequelID just for Sequels and as mentioned previously, it is also based on
    // release date rather than the timeline of the series as there are certain series such as fate where it is hard to
    // say which specific installment is prequel and which installment is sequel due to it sheer amount of references to
    // different series and the sheer number of installments that it currently has.
    Sequel *[]Sakuhin `json:"sequel_id"`
    // Parodies is almost certainly dedicated to doujin and is pretty self explanatory.
    Parodies  *[]SakuhinParodies `json:"parodies"`   // if nil then the series is an original series
    PageCount uint16             `json:"page_count"` // self explanatory
    // Seriel number of the installment release count starting from 1..n installments.
    InstallmentNumber uint32 `json:"installment_number"` // Is there any series with more than 4 billion installments?
    // This field was mostly for the frontend but I guess we can use it to see if the installment is a "special", "extra", etc.
    DisplayLabel db.DisplayLabel `json:"display_label"`
    ReleasedAt   time.Time       `json:"released_at"` // self explanatory
    UploadedAt   time.Time       `json:"uploaded_at"` // self explanatory
    UpdatedAt    time.Time       `json:"updated_at"`  // self explanatory

    // The following is not available in the table as a singleton and we need to fetch using foreign key and merge them.
    // The above instruction would be completed at the database layer so that we don't have to think about it.
    Representation SakuhinRepresentation `json:"representation"`
}

// Modern language don't have unions so I'm just gonna enforce in the frontend that a title must be provided, either En or Jp.
type StringWithLang struct {
    English  *string `json:"english"`
    Japanese *string `json:"japanese"`
}

type SakuhinTag struct {
    // TODO: Probably better to have a enum dedicated to themes and genres defining all possible combinations
    Themes  []string  `json:"themes"`  // Themes are stuff like "foot fetish", "tomboy", "gore", "ugly bastard", "slasher", etc.
    Genres  []string  `json:"genres"`  // Genres are stuff like "Horror", "Romance", "Yuri", "Isekai", "Comedy", "Action", etc.
    Circles *[]string `json:"circles"` // Circle is most likely gonna be nil for manga and list of circles/groups for doujin.
}

// Role agnostic strucutre which can be used to represent either an artist, author, etc.
// This is probably gonna be a table which is autogenerated whena new author/artist name is referenced while uploading a work.
// This way we can add autocomplete feature displaying the name when the user types while uploading new works in the frontend UI.
type SakuhinCreator struct {
    ID   uuid.UUID      `json:"id"`
    Name StringWithLang `json:"name"`
}

// Separate table
type SakuhinRepresentation struct {
    ChapterCount *uint64 `json:"chapter_count"` // will be nil in cases where we cannot count chapter such as in tankouban release, etc.
    Volumes      *[]struct {
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

// FIXME: If ACharacter is gonna be a table then we probably can define their role in their row itself and don't need this structure.
type SakuhinCharacter struct {
    Protagonists         []ACharacter `json:"protagonists"`
    SupportingCharacters []ACharacter `json:"supporting_characters"`
    SideCharacters       []ACharacter `json:"side_characters"`
    Antagonists          []ACharacter `json:"antagonists"`
    Unknown              []ACharacter `json:"unknown"`
}

// NOTE: THIS HAS TO BE A TABLE OF CHARACTER WITH EACH CHARACTER HAVING AN UNIQUE ID.
type ACharacter struct {
    EnglishTitle  *string      `json:"english_title"`  // if nil, then character name is probably not translated
    JapaneseTitle *string      `json:"japanese_title"` // if nil, then character is probably not of japanese origin
    Photos        *[]uuid.UUID `json:"photos"`         // if nil, then the character probably is a mythical faceless character or just not been revealed in the story yet
}

type SakuhinParodies struct {
    // TODO: Update this shit, I cannot guarentee that I would have the original series downloaded
    // that the doujin parodies (I'm not sure if parodies is the right word here, maybe parody-ed).
    OriginalID uuid.UUID    `json:"original_id"` // FIXME: Foreign key to the original series
    Characters []ACharacter `json:"characters"`
}

// Separate table
type SakuhinPublicationStat struct {
    SeriesID        uuid.UUID   `json:"series_id"` // foreign key
    InitialRelease  time.Time   `json:"initial_release"`
    LatestRelease   time.Time   `json:"latest_release"`
    CurrentStatus   db.Status   `json:"status"`
    ReleaseSchedule db.Schedule `json:"schedule"`
}

