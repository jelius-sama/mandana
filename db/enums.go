package db

import "github.com/jelius-sama/logger"

type SakuhinCategory uint32
type SakuhinDemographic uint32
type SakuhinContentRating uint32
type DisplayLabel uint8
type Status uint32
type Schedule uint32

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
const (
    StatusOngoing Status = iota
    StatusFinished
    StatusHiatus
    StatusAbandoned
    StatusUnknown
)

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

