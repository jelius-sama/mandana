package types

import (
    "time"

    "github.com/google/uuid"
    "github.com/jelius-sama/logger"
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
    CoverArts *[]uuid.UUID `json:"cover_arts"`
    // Tag is probably gonna be a new table with a foreign key to the row it represents, it contains a bunch of tags.
    Tags SakuhinTag `json:"tags"`
    // I'm not sure if category will come in any use but let's just keep it just in case to distinguish.
    Category SakuhinCategory `json:"category"`
    // Demographic is the target audience base of a specific work. Can be nil in some cases.
    Demographic   *SakuhinDemographic  `json:"demographic"`    // can be nil if unspecified demographic or multi-demographic
    ContentRating SakuhinContentRating `json:"content_rating"` // self explanatory
    Artist        []SakuhinCreator     `json:"artist"`         // self explanatory
    Author        []SakuhinCreator     `json:"author"`         // self explanatory
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
    PrequelID *[]uuid.UUID `json:"prequel_id"`
    // Same as previously defined field PrequelID just for Sequels and as mentioned previously, it is also based on
    // release date rather than the timeline of the series as there are certain series such as fate where it is hard to
    // say which specific installment is prequel and which installment is sequel due to it sheer amount of references to
    // different series and the sheer number of installments that it currently has.
    SequelID *[]uuid.UUID `json:"sequel_id"`
    // Parodies is almost certainly dedicated to doujin and is pretty self explanatory.
    Parodies  *[]SakuhinParodies `json:"parodies"`   // if nil then the series is an original series
    PageCount uint16             `json:"page_count"` // self explanatory
    // Seriel number of the installment release count starting from 1..n installments.
    InstallmentNumber uint32 `json:"installment_number"` // Is there any series with more than 4 billion installments?
    // This field was mostly for the frontend but I guess we can use it to see if the installment is a "special", "extra", etc.
    DisplayLabel DisplayLabel `json:"display_label"`
    ReleasedAt   time.Time    `json:"released_at"` // self explanatory
    UploadedAt   time.Time    `json:"uploaded_at"` // self explanatory
    UpdatedAt    time.Time    `json:"updated_at"`  // self explanatory
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

type SakuhinCategory uint32
type SakuhinDemographic uint32
type SakuhinContentRating uint32

const (
    CategoryDoujin SakuhinCategory = iota
    CategoryManga
)

const (
    DemographicShounen SakuhinDemographic = iota
    DemographicShoujo
    DemographicSeinen
    DemographicJosei
)

const (
    RatingSafe       SakuhinContentRating = iota // SFW Work
    RatingSuggestive                             // Not pornographic work but not SFW either
    RatingErotica                                // Soft pornographic work (blured censorship, holy light censorship, etc.)
    RatingNSFW                                   // Uncensored pornographic work
)

func (mc SakuhinCategory) String() string {
    switch mc {
    case CategoryDoujin:
        return "Doujin"
    case CategoryManga:
        return "Manga"
    }

    logger.Panic("unreachable")
    return ""
}

func (md SakuhinDemographic) String() string {
    switch md {
    case DemographicJosei:
        return "Josei"
    case DemographicSeinen:
        return "Seinen"
    case DemographicShoujo:
        return "Shoujo"
    case DemographicShounen:
        return "Shounen"
    }

    logger.Panic("unreachable")
    return ""
}

func (mcr SakuhinContentRating) String() string {
    switch mcr {
    case RatingErotica:
        return "Erotica"
    case RatingNSFW:
        return "Hentai"
    case RatingSafe:
        return "Safe"
    case RatingSuggestive:
        return "Suggestive"
    }

    logger.Panic("unreachable")
    return ""
}

// Separate table
type SakuhinRepresentation struct {
    SeriesID     uuid.UUID `json:"series_id"`     // foreign key
    ChapterCount *uint64   `json:"chapter_count"` // will be nil in cases where we cannot count chapter such as in tankouban release, etc.
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

type DisplayLabel uint8

const (
    LabelPart DisplayLabel = iota
    LabelOneShot
    LabelSpecial
    LabelExtra
    LabelSpinOff
    LabelVolume  // If tankouban release schedule
    LabelChapter // If individual chapters are released instead of more commonly used parts convention
    LabelCour    // Rarely used in Manga/Doujin context, but let's keep it for edge cases
    LabelUnknown
)

func (dl DisplayLabel) String() string {
    switch dl {
    case LabelCour:
        return "Cour"
    case LabelOneShot:
        return "One Shot"
    case LabelPart:
        return "Part"
    case LabelSpecial:
        return "Special"
    case LabelSpinOff:
        return "Spin Off"
    case LabelVolume:
        return "Volume"
    case LabelChapter:
        return "Chapter"
    case LabelExtra:
        return "Extra"
    case LabelUnknown:
        return "Unknown"
    }

    logger.Panic("unreachable")
    return ""
}

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
    SeriesID        uuid.UUID `json:"series_id"` // foreign key
    InitialRelease  time.Time `json:"initial_release"`
    LatestRelease   time.Time `json:"latest_release"`
    CurrentStatus   Status    `json:"status"`
    ReleaseSchedule Schedule  `json:"schedule"`
}

type Status uint32

const (
    StatusOngoing Status = iota
    StatusFinished
    StatusHiatus
    StatusAbandoned
    StatusUnknown
)

func (s Status) String() string {
    switch s {
    case StatusAbandoned:
        return "Abandoned"
    case StatusFinished:
        return "Finished"
    case StatusHiatus:
        return "Hiatus"
    case StatusOngoing:
        return "Ongoing"
    case StatusUnknown:
        return "Unknown"
    }

    logger.Panic("unreachable")
    return ""
}

type Schedule uint32

const (
    ScheduleWeekly Schedule = iota
    ScheduleBiWeekly
    ScheduleSemiWeekly
    ScheduleMonthly
    ScheduleBiMonthly
    ScheduleSemiMonthly
    ScheduleIrregular
    ScheduleTankoubon
    ScheduleOneShot
)

func (s Schedule) String() string {
    switch s {
    case ScheduleWeekly:
        return "Weekly"
    case ScheduleBiWeekly:
        return "Bi-Weekly"
    case ScheduleSemiWeekly:
        return "Semi-Weekly"
    case ScheduleSemiMonthly:
        return "Semi-Monthly"
    case ScheduleMonthly:
        return "Monthly"
    case ScheduleBiMonthly:
        return "Bi-Monthly"
    case ScheduleIrregular:
        return "Irregular"
    case ScheduleTankoubon:
        return "Tankoubon"
    case ScheduleOneShot:
        return "OneShot"
    }

    logger.Panic("unreachable")
    return ""
}

