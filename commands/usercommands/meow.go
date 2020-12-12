package usercommands

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Starshine113/covebotnt/crouter"
	"github.com/bwmarrin/discordgo"
)

var meows = []string{"350320411617263618", "350401541066194954", "350452775924269068", "350676297846095872", "350680941804388352", "350734766385004554", "350753356714606603", "350785335208181772", "351090256717807618", "351097460589854720", "352818252579078144", "352818268572221442", "358328186519420929", "361151997488922626", "364799899037859840", "375685680116006933", "380802572375818253", "415956883568984065", "447075093944860682", "496818729829531648", "585552852886356076", "585786840309694475", "587999110959988739", "587999110997606420", "587999111005995018", "587999111027097601", "587999111081361417", "587999111161315361", "587999111265910825", "587999111307984909", "587999111358185476", "587999111362379786", "587999111366836235", "587999111370899466", "587999111404322816", "587999111442071552", "587999111551385641", "587999112113422336", "587999112688041994", "599204881122852865", "600787523421208577", "621634488220385290", "623818335041224704", "645414359337664516", "653601602732818442", "778573046134865970", "778613156523278336", "380069533752360970", "380069565016834058", "380070162092654596", "380070955302780928", "396448565213396992", "398175245934133269", "396451934921555978", "406922357312192553", "407202980610310154", "380089286831374337", "380110519215980559", "380404639880839179", "409770961831854080", "411237179231174667", "496335045007638528", "518800392331722773", "390962614840328192", "391358218971906050", "411236858538754050", "411239608651874315", "411254866049237003", "411644195359424515", "415588670331027478", "417357486329298954", "424575186814369792", "424585599685623809", "426067749756862495", "427470977895759883", "427472210383601675", "442778936452055040", "445892579763159060", "446287918063943687", "446307019633459210", "447074986516021258", "454580034234089472", "463767725655719956", "498772285918937088", "498772285990240257"}

func meow(ctx *crouter.Ctx) (err error) {
	r := meows[rand.Intn(len(meows)-1)]

	url := fmt.Sprintf("https://cdn.discordapp.com/emojis/%v.png", r)
	resp, err := http.Get(url)
	if err != nil {
		_, err = ctx.CommandError(err)
		return err
	}
	defer resp.Body.Close()

	file := &discordgo.File{
		Name:   fmt.Sprintf("%v.png", r),
		Reader: resp.Body,
	}

	_, err = ctx.Send(&discordgo.MessageSend{Files: []*discordgo.File{file}})
	return err
}
