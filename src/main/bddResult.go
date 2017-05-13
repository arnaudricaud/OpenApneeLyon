package main
	
	import (
	"strconv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"bufio"
	"strings"
	"time"
	)
	
	func (base Bdd) displayClassement(){
		base.resultat, base.err = base.db.Query("SELECT * FROM classement")
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(base.err)
		}
		defer base.resultat.Close()
		
		var info [11]string

		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9], &info[10])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
		fmt.Println(info[0] + "|" + info[1]+ "|" + info[2]+ "|" + info[3] + "|" + info[4]+ "|" + info[5]+ "|" + info[6]+ "|" + info[7]+ "|" + info[8]+ "|" + info[9]+ "|" + info[10])
		}
	}
	
	/*
	* 		Bdd.searchCompetiteur:
	* Paramètres:
	*	- col_num: 	numéro de la colonne sur laquelle effectuer la recherche (ex: 2 => prénom).
	*	- value:	valeur à rechercher dans la colonne choisie.
	*
	* Description: 		
	*		Méthode permettant de rechercher un compétiteur en
	* 		competiteurs de la base de données
	*/
	
	func (base Bdd) searchCompetiteurClassement(col_num int, value string){
		
		var id_col string
		var searchValue string
		
		searchValue = fmt.Sprint("'%",value,"%'")
		id_col, value = col_id2name(col_num, value)
		
		base.resultat, base.err = base.db.Query(fmt.Sprint("SELECT * FROM classement WHERE ", id_col, " LIKE ", searchValue))
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(base.err)
		}
		defer base.resultat.Close()
		
		var info [11]string

		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9], &info[10])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
		fmt.Println(info[0] + "|" + info[1]+ "|" + info[2]+ "|" + info[3] + "|" + info[4]+ "|" + info[5]+ "|" + info[6]+ "|" + info[7]+ "|" + info[8]+ "|" + info[9]+ "|" + info[10])
		}
	}
	
	
		/*
	* 		Bdd.addCompetiteur:
	* Paramètres:
	*	- comp: 	Les informations du compétiteur à ajouter sous la
	*				forme d'une structure de type "competiteur"
	*
	* Description: 		
	*		Méthode permettant d'ajouter un compétiteur dans la 
	* 		base de données
	*/

	func (base Bdd) addCompetiteurClassement(board *Classement){
		
		_, base.err = base.db.Exec("INSERT INTO classement ( id, prenom, nom, sexe, equipe, epreuve, annonce, resultat, place, rslt, plc) VALUES(" +
		strconv.Itoa(board.id) + ",'" +
		board.prenom + "','" +
		board.nom + "','" +
		board.sexe + "','" +
		board.equipe + "','" +
		board.epreuve + "'," +
		strconv.Itoa(board.annonce) + "," +
		strconv.Itoa(board.resultat) + "," +
		strconv.Itoa(board.place) + "," +
		strconv.Itoa(board.rslt) + "," +
		strconv.Itoa(board.plc) + ")")
	
		if base.err != nil {
			fmt.Println("Echec lors de l'ajout1 : "+ board.nom +" "+ board.prenom, base.err)
			} else {
			fmt.Println("Ajout validé du resulat compétiteur " + board.nom +" "+ board.prenom)
		}
	}
	
	/*
	* 		Bdd.deleteCompetiteur:
	* Paramètres:
	*	- col_num: 	numéro de la colonne sur laquelle effectuer la recherche (ex: 2 => prénom).
	*	- value:	valeur à rechercher dans la colonne choisie.
	*
	* Description: 		
	*		Méthode permettant de supprimer les compétiteurs en fonction des critères
	*		en entrée.
	*/

	func (base Bdd) deleteCompetiteurClassement(col_num int, value string){
		var id_col string
		value = fmt.Sprint("'",value,"'")
		
		if col_num==1 {
			id_col = "id"		
		} else if col_num==2{
			id_col = "equipe"		
		}
		
		if !(col_num < 1 && col_num > 2){
			_, base.err = base.db.Exec("DELETE FROM classement WHERE " + id_col + " = " + value)
			if base.err != nil {
				fmt.Println("Echec lors de la suppression: \n", base.err)
				} else {
				fmt.Println("Suppression des competiteurs avec " + id_col + " = " + value)
			}
		} else {
			err := "Le numéro entré est invalide!"
			fmt.Println(err);
		}
	}
	
	
	/*
	* 		Bdd.importCompetiteur:
	* Paramètres:
	*	- chemin: 	Chemin du fichier à importer avec le nom du fichier et l'extension.
	*
	* Description: 		
	*		Méthode permettant d'importer les compétiteurs contenu dans un fichier CSV
	*/
	
	func (base Bdd) importResultat(){
		file, err := os.Open("import/classement.csv")
		if err != nil {
			println("Impossible d'ouvrir le fichier \"classement.csv\" dans le dossier import")
			log.Fatal(err)
		}
		defer file.Close()
	
		var firstCall bool	
		firstCall = true
		var res int
		var place int
		var plc int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			info := strings.Split(scanner.Text(), ";")
			if !firstCall{
			temps,er := strconv.Atoi(info[6])
			idd,errr := strconv.Atoi(info[0])
			annonce := base.recupAnnonce(info[1],info[2],info[3],info[5])
			if er != nil {
			log.Fatal(er)
			}
			if errr != nil {
			log.Fatal(errr)
			}
			switch(info[5]){
			case "spd": 
			res=calculResultat("spd",annonce,info[6])
			break
			case "1650":
			res=calculResultat("1650",annonce,info[6])
			break
			case "dnf":
			res=calculResultat("dnf",annonce,info[6])
			break
			case "dwf":
			res=calculResultat("dwf",annonce,info[6])
			break
			case "sta":
			res=calculResultat("sta",annonce,info[6])
			break
			}
				
			classemt := newClassement(idd, info[1], info[2], info[3], info[4], info[5],annonce, temps,place,res,plc)
			base.addCompetiteurClassement(classemt)
			}
			firstCall = false
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}	
	}
	/*
	*
	*
	*
	*/
	func (base Bdd)recupAnnonce(prenom string, nom string, sexe string, epreuve string)(int){
		var id_col string
		id_col, prenom = col_id2name2(2, prenom)
		var id_col2 string
		id_col2, nom = col_id2name2(3, nom)
		var id_col3 string
		id_col3, sexe = col_id2name2(4, sexe)
		base.resultat, base.err = base.db.Query("SELECT * FROM competiteurs WHERE " + id_col + " = " + prenom + " AND " + id_col2 + " = " + nom + " AND " + id_col3 + " = " + sexe)
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(base.err)
		}
		defer base.resultat.Close()
		
		var info [10]string
		var resultat int
		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
			if (epreuve==info[6]){
			resultat,_ = strconv.Atoi(info[7])
			}else if (epreuve==info[8]){
			resultat,_ = strconv.Atoi(info[9])
			} else{
			resultat = 0
			}	
		}
		return resultat
	}
		
	/*
	* 		col_id2name:
	* Paramètres:
	*	- col_num:  Numéro de la colonne sur laquelle effectuer la modification (ex: 2 => prénom).
	*	- value:	Nouvelle valeur à entrée pour la colonne choisie.
	*
	* Description: 		
	*		Méthode permettant à partir d'un numéro de colonne, de retourner le nom de la colonne.
	*		De plus, la valeur entrée ("value") est retournée au format adéquat pour une requête SQL
	*		(Ex: "VariableString" => "'VariableString'")
	*/
	
	
	func col_id2name2(col_num int, value string)(string, string){
		var col_idr string
		
		switch col_num{
		    case 1:
				col_idr = "id"
				value = fmt.Sprint("'",value,"'")
			case 2:
				col_idr = "prenom"
				value = fmt.Sprint("'",value,"'")
			case 3:
				col_idr = "nom"
				value = fmt.Sprint("'",value,"'")
			case 4:
				col_idr = "sexe"
				value = fmt.Sprint("'",value,"'")
			case 5:
				col_idr = "equipe"
				value = fmt.Sprint("'",value,"'")
			case 6:
				col_idr = "epreuve"
				value = fmt.Sprint("'",value,"'")
			case 7:
				col_idr = "annonce"
				value = fmt.Sprint("'",value,"'")
			case 8:
				col_idr = "resultat"
				value = fmt.Sprint("'",value,"'")
			case 9:
				col_idr = "place"
				value = fmt.Sprint("'",value,"'")
			case 10:
				col_idr = "rslt"
				value = fmt.Sprint("'",value,"'")
			case 11:
				col_idr = "plc"
				value = fmt.Sprint("'",value,"'")
			default:
				log.Fatal("Numéro invalide")
			}
		return col_idr, value
	}
	
	
	/*
	*
	*
	*
	*/
	func (base Bdd) exportClassement(value string){ 
		t := time.Now()
			date := fmt.Sprint(t.Year(),"_",int(t.Month()),"_", t.Day(),"_",t.Hour(),"_", t.Minute(),"_", t.Second())
		
		file, err := os.Create(fmt.Sprint("export/archives/",date,"-",value,".csv"))
			if err != nil {
				fmt.Println("Erreur lors de la création du fichier. Avez vous créé un dossier \"export/archives\" dans le dossier de l'application?")
				log.Fatal(err)
			}
		file2, err := os.Create(fmt.Sprint("export/Classement-",value,".csv"))
		if err != nil {
			fmt.Println("Erreur lors de la création du fichier. Avez vous créé un dossier \"export\" dans le dossier de l'application?")
			log.Fatal(err)
		}
			
			
			//calcul de la place equipe
			switch(value){
			case "spd": 
			base.calculPlace("spd")
			break
			case "1650":
			base.calculPlace("1650")
			break
			case "dnf":
			base.calculPlace("dnf")
			break
			case "dwf":
			base.calculPlace("dwf")
			break
			case "sta":
			base.calculPlace("sta")
			break
			}
			
		var id_col string 
		id_col, value = col_id2name2(6, value)
	base.resultat, base.err = base.db.Query(fmt.Sprint("SELECT * FROM classement WHERE ", id_col, " = ", value," ORDER BY sexe ASC, resultat DESC"))
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête 1")
		}
		defer base.resultat.Close()
	var info [11]string
	var numPlaceF int =1
	var numPlaceH int =1
	var sexe string ="F" 
	
	
	file.WriteString(fmt.Sprint("\xEF\xBB\xBFId; Prenom; Nom; Sexe; Equipe; Epreuve; Annonce; Resultat; Place; Resultat pris en compte equipe; Place Equipe\r\n"))
	file2.WriteString(fmt.Sprint("\xEF\xBB\xBFId; Prenom; Nom; Sexe; Equipe; Epreuve; Annonce; Resultat; Place; Resultat pris en compte equipe; Place Equipe\r\n"))	
		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9], &info[10])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)}				
			
			if(info[3]==sexe){
				info[8]=strconv.Itoa(numPlaceF)
				numPlaceF=numPlaceF+1
			}else{
				info[8]=strconv.Itoa(numPlaceH)
				numPlaceH=numPlaceH+1}
				//id,_:=strconv.Atoi(info[0])
			//base.modifResult(id,9,info[8])

		file.WriteString(fmt.Sprint(info[0],";",info[1],";", info[2],";", info[3],";", info[4],";", info[5],";", info[6],";", info[7],";", info[8],";", info[9],";", info[10],"\r\n"))
		file2.WriteString(fmt.Sprint(info[0],";",info[1],";", info[2],";", info[3],";", info[4],";", info[5],";", info[6],";", info[7],";", info[8],";", info[9],";", info[10],"\r\n"))
		}
	}
	
	func (base Bdd) modifResult(id_comp int, col_num int, newvalue string){
		col_id, value := col_id2name2(col_num, newvalue)
		id := strconv.Itoa(id_comp)
		_, base.err = base.db.Exec("UPDATE classement SET "  + col_id + " = " + value +  " WHERE id = " + id)
	
		if base.err != nil {
			fmt.Println("Echec lors de l'ajout : ", base.err)
			} else {
			//fmt.Println("Modification du competiteur " + strconv.Itoa(id_comp) + " avec " + col_id + " = " + value)
		}
	}
	
	/*
	*
	*/
	func (base Bdd) calculPlace(epreuve string){
	var id_col string 
		id_col, epreuve = col_id2name2(6, epreuve)
	base.resultat, base.err = base.db.Query(fmt.Sprint("SELECT * FROM classement WHERE ", id_col, " = ", epreuve," ORDER BY sexe ASC, rslt DESC"))
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête 1")
		}
	var info [11]string
	var numPlaceF int =1
	var numPlaceH int =1
	var sexe string ="F" 
	var tabClassement []*Classement
	var nextResult *Classement
		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9], &info[10])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)}				
				
		if(info[3]==sexe){
				info[10]=strconv.Itoa(numPlaceF)
				numPlaceF=numPlaceF+1
				
			}else{
				info[10]=strconv.Itoa(numPlaceH)
				numPlaceH=numPlaceH+1
				}
				
				
				id,_:=strconv.Atoi(info[0])
				annonce,_ := strconv.Atoi(info[6])
				resultat,_ := strconv.Atoi(info[7])
				place,_ := strconv.Atoi(info[8])
				rslt,_ := strconv.Atoi(info[9])
				plc,_ := strconv.Atoi(info[10])
				
				nextResult = newClassement(id, info[1], info[2], info[3], info[4], info[5], annonce, resultat, place, rslt, plc)
				tabClassement = append(tabClassement, nextResult)
		}
		base.resultat.Close()
		 
			for i := 0; i < len(tabClassement); i++{
				base.modifResult(tabClassement[i].id ,11,strconv.Itoa(tabClassement[i].plc))
			}

	}
	
	
	func calculResultat(epreuve string, annonce int, resultat string)(int){
	var sMin int =0 
	var sMax int =0
	var res int
	var result int
	var tot int
	var tot2 int
	
	var tab[] *ConfigurationEpreuve
	result,_ =strconv.Atoi(resultat)
	tab=getConfigurationEpreuve1()
	
	for i := 0; i < 5; i++{
	if (tab[i].id==epreuve){
	sMin=tab[i].seuilMin
	sMax=tab[i].seuilMax
	}
	
	max:=annonce+sMax
	min:=annonce+sMin
	if(result>max){
	switch(epreuve){
	case "spd": 
	tot =(result-(annonce+20))*3
	break
	case "1650":
	tot = (result-(annonce+60))*3
	break
	case "dnf":
	tot = (annonce+25)
	break
	case "dwf":
	tot = (annonce+25)
	break
	case "sta":
	tot = (annonce+60)
	break	
	}
	res=tot
	}else if(result<min){
	switch(epreuve){
	case "spd": 
	tot2=annonce-10
	break
	case "1650":
	tot2=annonce-30
	break
	case "dnf":
	tot2=((annonce-25)-result)*3
	break
	case "dwf":
	tot2=((annonce-25)-result)*3
	break
	case "sta":
	tot2=((annonce-60)-result)*3
	break	
	}
	res=tot2
	}else if (result == annonce){
	res= result
	}
	
	}
	
	return res
	}
	
	/*
	*
	*
	*/
	func getConfigurationEpreuve1()([]*ConfigurationEpreuve){
	file, err := os.Open("config/ConfigurationEpreuve.csv")
	if err != nil {
		fmt.Println("Impossible d'ouvrir le fichier \"ConfigurationEpreuve\": " )
		log.Fatal(err)
	}
	defer file.Close()
	
	
	var firstCall bool
	firstCall = true
	var nextconfig *ConfigurationEpreuve
	
	scanner := bufio.NewScanner(file)
	//On clear l'ancien tableau:
	var tabConfig[] *ConfigurationEpreuve 
	tabConfig=tabConfig[:0]
	
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), ";")
		if !firstCall{
		ordre, _ := strconv.Atoi(info[0])
		seuilMin, _ := strconv.Atoi(info[2])
		seuilMax, _ := strconv.Atoi(info[3])
		nbParPassage, _ := strconv.Atoi(info[4])
		dureeEchauffement, _ := strconv.Atoi(info[5])
		dureePassage, _ := strconv.Atoi(info[6])
		dureeAppel, _ := strconv.Atoi(info[7])
		surveillance, _ := strconv.Atoi(info[8])
		battementSerie, _ := strconv.Atoi(info[9])
		battementEpreuve, _ := strconv.Atoi(info[10])
	
		nextconfig = newConfigurationEpreuve(ordre, info[1], seuilMin, seuilMax, nbParPassage, dureeEchauffement, dureePassage, dureeAppel, surveillance,
												battementSerie,battementEpreuve, info[11])
		tabConfig=append(tabConfig,nextconfig)
		}
		firstCall = false
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tabConfig
}

	
	