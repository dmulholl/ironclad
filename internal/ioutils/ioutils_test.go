package ioutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var deBelloGallico = `
Gallia est omnis divisa in partes tres, quarum unam incolunt Belgae, aliam Aquitani, tertiam qui
ipsorum lingua Celtae, nostra Galli appellantur. Hi omnes lingua, institutis, legibus inter se
differunt. Gallos ab Aquitanis Garumna flumen, a Belgis Matrona et Sequana dividit. Horum omnium
fortissimi sunt Belgae, propterea quod a cultu atque humanitate provinciae longissime absunt,
minimeque ad eos mercatores saepe commeant atque ea quae ad effeminandos animos pertinent important,
proximique sunt Germanis, qui trans Rhenum incolunt, quibuscum continenter bellum gerunt. Qua de
causa Helvetii quoque reliquos Gallos virtute praecedunt, quod fere cotidianis proeliis cum Germanis
contendunt, cum aut suis finibus eos prohibent aut ipsi in eorum finibus bellum gerunt. Eorum una,
pars, quam Gallos obtinere dictum est, initium capit a flumine Rhodano, continetur Garumna flumine,
Oceano, finibus Belgarum, attingit etiam ab Sequanis et Helvetiis flumen Rhenum, vergit ad
septentriones. Belgae ab extremis Galliae finibus oriuntur, pertinent ad inferiorem partem fluminis
Rheni, spectant in septentrionem et orientem solem. Aquitania a Garumna flumine ad Pyrenaeos montes
et eam partem Oceani quae est ad Hispaniam pertinet; spectat inter occasum solis et septentriones.

Apud Helvetios longe nobilissimus fuit et ditissimus Orgetorix. Is M. Messala, et P. M. Pisone
consulibus regni cupiditate inductus coniurationem nobilitatis fecit et civitati persuasit ut de
finibus suis cum omnibus copiis exirent: perfacile esse, cum virtute omnibus praestarent, totius
Galliae imperio potiri. Id hoc facilius iis persuasit, quod undique loci natura Helvetii
continentur: una ex parte flumine Rheno latissimo atque altissimo, qui agrum Helvetium a Germanis
dividit; altera ex parte monte Iura altissimo, qui est inter Sequanos et Helvetios; tertia lacu
Lemanno et flumine Rhodano, qui provinciam nostram ab Helvetiis dividit. His rebus fiebat ut et
minus late vagarentur et minus facile finitimis bellum inferre possent; qua ex parte homines
bellandi cupidi magno dolore adficiebantur. Pro multitudine autem hominum et pro gloria belli atque
fortitudinis angustos se fines habere arbitrabantur, qui in longitudinem milia passuum CCXL, in
latitudinem CLXXX patebant.

His rebus adducti et auctoritate Orgetorigis permoti constituerunt ea quae ad proficiscendum
pertinerent comparare, iumentorum et carrorum quam maximum numerum coemere, sementes quam maximas
facere, ut in itinere copia frumenti suppeteret, cum proximis civitatibus pacem et amicitiam
confirmare. Ad eas res conficiendas biennium sibi satis esse duxerunt; in tertium annum profectionem
lege confirmant. Ad eas res conficiendas Orgetorix deligitur. Is sibi legationem ad civitates
suscipit. In eo itinere persuadet Castico, Catamantaloedis filio, Sequano, cuius pater regnum in
Sequanis multos annos obtinuerat et a senatu populi Romani amicus appellatus erat, ut regnum in
civitate sua occuparet, quod pater ante habuerit; itemque Dumnorigi Haeduo, fratri Diviciaci, qui eo
tempore principatum in civitate obtinebat ac maxime plebi acceptus erat, ut idem conaretur persuadet
eique filiam suam in matrimonium dat. Perfacile factu esse illis probat conata perficere, propterea
quod ipse suae civitatis imperium obtenturus esset: non esse dubium quin totius Galliae plurimum
Helvetii possent; se suis copiis suoque exercitu illis regna conciliaturum confirmat. Hac oratione
adducti inter se fidem et ius iurandum dant et regno occupato per tres potentissimos ac firmissimos
populos totius Galliae sese potiri posse sperant.

Ea res est Helvetiis per indicium enuntiata. Moribus suis Orgetoricem ex vinculis causam dicere
coegerunt; damnatum poenam sequi oportebat, ut igni cremaretur. Die constituta causae dictionis
Orgetorix ad iudicium omnem suam familiam, ad hominum milia decem, undique coegit, et omnes clientes
obaeratosque suos, quorum magnum numerum habebat, eodem conduxit; per eos ne causam diceret se
eripuit. Cum civitas ob eam rem incitata armis ius suum exequi conaretur multitudinemque hominum ex
agris magistratus cogerent, Orgetorix mortuus est; neque abest suspicio, ut Helvetii arbitrantur,
quin ipse sibi mortem consciverit.

Post eius mortem nihilo minus Helvetii id quod constituerant facere conantur, ut e finibus suis
exeant. Ubi iam se ad eam rem paratos esse arbitrati sunt, oppida sua omnia, numero ad duodecim,
vicos ad quadringentos, reliqua privata aedificia incendunt; frumentum omne, praeter quod secum
portaturi erant, comburunt, ut domum reditionis spe sublata paratiores ad omnia pericula subeunda
essent; trium mensum molita cibaria sibi quemque domo efferre iubent. Persuadent Rauracis et
Tulingis et Latobrigis finitimis, uti eodem usi consilio oppidis suis vicisque exustis una cum iis
proficiscantur, Boiosque, qui trans Rhenum incoluerant et in agrum Noricum transierant Noreiamque
oppugnabant, receptos ad se socios sibi adsciscunt.

Erant omnino itinera duo, quibus itineribus domo exire possent: unum per Sequanos, angustum et
difficile, inter montem Iuram et flumen Rhodanum, vix qua singuli carri ducerentur, mons autem
altissimus impendebat, ut facile perpauci prohibere possent; alterum per provinciam nostram, multo
facilius atque expeditius, propterea quod inter fines Helvetiorum et Allobrogum, qui nuper pacati
erant, Rhodanus fluit isque non nullis locis vado transitur. Extremum oppidum Allobrogum est
proximumque Helvetiorum finibus Genava. Ex eo oppido pons ad Helvetios pertinet. Allobrogibus sese
vel persuasuros, quod nondum bono animo in populum Romanum viderentur, existimabant vel vi coacturos
ut per suos fines eos ire paterentur. Omnibus rebus ad profectionem comparatis diem dicunt, qua die
ad ripam Rhodani omnes conveniant.

Caesari cum id nuntiatum esset, eos per provinciam nostram iter facere conari, maturat ab urbe
proficisci et quam maximis potest itineribus in Galliam ulteriorem contendit et ad Genavam pervenit.
Provinciae toti quam maximum potest militum numerum imperat erat omnino in Gallia ulteriore legio
una, pontem, qui erat ad Genavam, iubet rescindi. Ubi de eius adventu Helvetii certiores facti sunt,
legatos ad eum mittunt nobilissimos civitatis, cuius legationis Nammeius et Verucloetius principem
locum obtinebant, qui dicerent sibi esse in animo sine ullo maleficio iter per provinciam facere,
propterea quod aliud iter haberent nullum: rogare ut eius voluntate id sibi facere liceat. Caesar,
quod memoria tenebat L. Cassium consulem occisum exercitumque eius ab Helvetiis pulsum et sub iugum
missum, concedendum non putabat; neque homines inimico animo, data facultate per provinciam itineris
faciundi, temperaturos ab iniuria et maleficio existimabat. Tamen, ut spatium intercedere posset dum
milites quos imperaverat convenirent, legatis respondit diem se ad deliberandum sumpturum.

Interea ea legione quam secum habebat militibusque, qui ex provincia convenerant, a lacu Lemanno,
qui in flumen Rhodanum influit, ad montem Iuram, qui fines Sequanorum ab Helvetiis dividit, milia
passuum decem novem murum in altitudinem pedum sedecim fossamque perducit. Eo opere perfecto
praesidia disponit, castella communit, quo facilius, si se invito transire conentur, prohibere
possit. Ubi ea dies quam constituerat cum legatis venit et legati ad eum reverterunt, negat se more
et exemplo populi Romani posse iter ulli per provinciam dare et, si vim facere conentur,
prohibiturum ostendit. Helvetii ea spe deiecti navibus iunctis ratibusque compluribus factis, alii
vadis Rhodani, qua minima altitudo fluminis erat, non numquam interdiu, saepius noctu si perrumpere
possent conati, operis munitione et militum concursu et telis repulsi, hoc conatu destiterunt.

Relinquebatur una per Sequanos via, qua Sequanis invitis propter angustias ire non poterant. His cum
sua sponte persuadere non possent, legatos ad Dumnorigem Haeduum mittunt, ut eo deprecatore a
Sequanis impetrarent. Dumnorix gratia et largitione apud Sequanos plurimum poterat et Helvetiis erat
amicus, quod ex ea civitate Orgetorigis filiam in matrimonium duxerat, et cupiditate regni adductus
novis rebus studebat et quam plurimas civitates suo beneficio habere obstrictas volebat. Itaque rem
suscipit et a Sequanis impetrat ut per fines suos Helvetios ire patiantur, obsidesque uti inter sese
dent perficit: Sequani, ne itinere Helvetios prohibeant, Helvetii, ut sine maleficio et iniuria
transeant.

Caesari renuntiatur Helvetiis esse in animo per agrum Sequanorum et Haeduorum iter in Santonum fines
facere, qui non longe a Tolosatium finibus absunt, quae civitas est in provincia. Id si fieret,
intellegebat magno cum periculo provinciae futurum ut homines bellicosos, populi Romani inimicos,
locis patentibus maximeque frumentariis finitimos haberet. Ob eas causas ei munitioni quam fecerat
T. Labienum legatum praeficit; ipse in Italiam magnis itineribus contendit duasque ibi legiones
conscribit et tres, quae circum Aquileiam hiemabant, ex hibernis educit et, qua proximum iter in
ulteriorem Galliam per Alpes erat, cum his quinque legionibus ire contendit. Ibi Ceutrones et
Graioceli et Caturiges locis superioribus occupatis itinere exercitum prohibere conantur.
Compluribus his proeliis pulsis ab Ocelo, quod est oppidum citerioris provinciae extremum, in fines
Vocontiorum ulterioris provinciae die septimo pervenit; inde in Allobrogum fines, ab Allobrogibus in
Segusiavos exercitum ducit. Hi sunt extra provinciam trans Rhodanum primi.
`

func TestRoundtripFileEmptyInput(t *testing.T) {
	file, err := os.CreateTemp("", "ironclad-fileio-test")
	require.NoError(t, err, "failed to create temp fileio test file")
	defer os.Remove(file.Name())

	err = file.Close()
	require.NoError(t, err, "failed to close temp fileio test file")

	plaintext := []byte{}

	err = Save(file.Name(), "abc123", plaintext)
	require.NoError(t, err, "failed to save fileio test file")

	result, err := Load(file.Name(), "abc123")
	require.NoError(t, err, "failed to load fileio test file")

	require.Equal(t, plaintext, result)
}

func TestRoundtripFileShortInput(t *testing.T) {
	file, err := os.CreateTemp("", "ironclad-fileio-test")
	require.NoError(t, err, "failed to create temp fileio test file")
	defer os.Remove(file.Name())

	err = file.Close()
	require.NoError(t, err, "failed to close temp fileio test file")

	plaintext := []byte("foo bar baz")

	err = Save(file.Name(), "abc123", plaintext)
	require.NoError(t, err, "failed to save fileio test file")

	result, err := Load(file.Name(), "abc123")
	require.NoError(t, err, "failed to load fileio test file")

	require.Equal(t, plaintext, result)
}

func TestRoundtripFileLongInput(t *testing.T) {
	file, err := os.CreateTemp("", "ironclad-fileio-test")
	require.NoError(t, err, "failed to create temp fileio test file")
	defer os.Remove(file.Name())

	err = file.Close()
	require.NoError(t, err, "failed to close temp fileio test file")

	plaintext := []byte(deBelloGallico)

	err = Save(file.Name(), "abc123", plaintext)
	require.NoError(t, err, "failed to save fileio test file")

	result, err := Load(file.Name(), "abc123")
	require.NoError(t, err, "failed to load fileio test file")

	require.Equal(t, plaintext, result)
}

func TestRoundtripEmptyInput(t *testing.T) {
	plaintext := []byte{}

	data, err := Encrypt("abc123", plaintext)
	require.NoError(t, err, "failed to encrypt data")

	result, err := Decrypt("abc123", data)
	require.NoError(t, err, "failed to decrypt data")

	require.Equal(t, plaintext, result)
}

func TestRoundtripShortInput(t *testing.T) {
	plaintext := []byte("foo bar baz")

	data, err := Encrypt("abc123", plaintext)
	require.NoError(t, err, "failed to encrypt data")

	result, err := Decrypt("abc123", data)
	require.NoError(t, err, "failed to decrypt data")

	require.Equal(t, plaintext, result)
}

func TestRoundtripLongInput(t *testing.T) {
	plaintext := []byte(deBelloGallico)

	data, err := Encrypt("abc123", plaintext)
	require.NoError(t, err, "failed to encrypt data")

	result, err := Decrypt("abc123", data)
	require.NoError(t, err, "failed to decrypt data")

	require.Equal(t, plaintext, result)
}
