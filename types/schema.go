package types

import (
    "time"

    "github.com/google/uuid"
    "github.com/jelius-sama/logger"
)

// Sakuhin is designed with Manga and Doujin in mind, it can probably represent other work such as LN too.
type Sakuhin struct {
    ID                uuid.UUID            `json:"id"`
    Title             StringWithLang       `json:"title"`
    Description       StringWithLang       `json:"description"`
    CoverArts         []string             `json:"cover_arts"`
    Tags              SakuhinTag           `json:"tags"`
    Category          SakuhinCategory      `json:"category"`
    Demographic       *SakuhinDemographic  `json:"demographic"` // can be nil if unspecified demographic or multi-demographic
    ContentRating     SakuhinContentRating `json:"content_rating"`
    Artist            []string             `json:"artist"`
    Author            []string             `json:"author"`
    Character         SakuhinCharacter     `json:"character"`
    Language          string               `json:"language"`
    SeriesID          uuid.UUID            `json:"series_id"` // All installments/parts of a series share this same ID
    PrequelID         *[]uuid.UUID         `json:"prequel_id"`
    SequelID          *[]uuid.UUID         `json:"sequel_id"`
    Parodies          *[]SakuhinParodies   `json:"parodies"` // if nil then original series, this is mostly applicable to doujin
    PageCount         uint64               `json:"page_count"`
    InstallmentNumber uint64               `json:"installment_number"`
    DisplayLabel      DisplayLabel         `json:"display_label"`
    ReleasedAt        time.Time            `json:"released_at"`
    UploadedAt        time.Time            `json:"uploaded_at"`
    UpdatedAt         time.Time            `json:"updated_at"`
}

type StringWithLang struct {
    English  string `json:"english"`
    Japanese string `json:"japanese"`
}

type SakuhinTag struct {
    // TODO: Probably better to have a enum dedicated to themes and genres defining all possible combinations
    Themes  []string  `json:"themes"`
    Genres  []string  `json:"genres"`
    Circles *[]string `json:"circles"` // Circle is nil for manga and list of circles/groups for doujin
}

type SakuhinCategory uint32
type SakuhinDemographic uint32
type SakuhinContentRating uint32

const (
    CategoryDoujin SakuhinCategory = iota
    CategorySakuhin
)

const (
    DemographicShounen SakuhinDemographic = iota
    DemographicShoujo
    DemographicSeinen
    DemographicJosei
)

const (
    RatingSafe       SakuhinContentRating = iota
    RatingErotica                         // Soft pornographic work (blured censorship, holy light censorship, etc.)
    RatingSuggestive                      // Not pornographic work
    RatingNSFW                            // Uncensored pornographic work
)

func (mc SakuhinCategory) String() string {
    switch mc {
    case CategoryDoujin:
        return "Doujin"
    case CategorySakuhin:
        return "Sakuhin"
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
}

type ACharacter struct {
    EnglishTitle  *string   `json:"english_title"`  // if nil, then character name is probably not translated
    JapaneseTitle *string   `json:"japanese_title"` // if nil, then character is probably not of japanese origin
    Photos        *[]string `json:"photos"`         // if nil, then the character probably is a mythical faceless character or just not been revealed in the story yet
}

type SakuhinParodies struct {
    OriginalID uuid.UUID    `json:"original_id"`
    Characters []ACharacter `json:"characters"`
}

// Separate table
type SakuhinPublicationStat struct {
    SeriesID        uuid.UUID `json:"series_id"` // foreign key
    InitialRelease  time.Time `json:"initial_release"`
    LatestRelease   time.Time `json:"latest_release"`
    CurrentStatus   Status    `json:"status"`
    ReleaseSchedule Schedule  `json:"Schedule"`
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

