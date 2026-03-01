package constants

/*
// Since golang doesn't support struct to be constant or immutable
// we will use C to make it constant.

typedef struct {
    char *English;
    char *Japanese;
} AppNameT;

const AppNameT AppName = {
    .English = "Mandana",
    .Japanese = "漫棚"
};
*/
import "C"

type AppNameT struct {
    English  string
    Japanese string
}

const (
    Version string = "1.0.0"
)

func AppName() AppNameT {
    return AppNameT{
        English:  C.GoString((*C.char)(C.AppName.English)),
        Japanese: C.GoString((*C.char)(C.AppName.Japanese)),
    }
}

