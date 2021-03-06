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
"regexp"
)
	
	

	/*
	* 		Bdd.displayCompetiteur:
	* Description: 	
	* 		Méthode permettant d'afficher l'integralité des
	* 		compétiteurs contenus dans la table "competiteurs".
	*/
	
	func (base Bdd) displayCompetiteur(){
		
		//REQUÊTE
		base.resultat, base.err = base.db.Query("SELECT * FROM competiteurs")
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(base.err)
		}
		defer base.resultat.Close()
		
		var info [10]string
		// RÉCUPÉRATION DES RESULTATS
		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
			
		//AFFICHAGE
		fmt.Println(info[0] + "|" + info[1]+ "|" + info[2]+ "|" + info[3] + "|" + info[4]+ "|" + info[5]+ "|" + info[6]+ "|" + info[7]+ "|" + info[8]+ "|" + info[9])
		}
	}
	
	/*
	* 		Bdd.searchCompetiteur:
	* Paramètres:
	*	- col_num: 	numéro de la colonne sur laquelle effectuer la recherche (ex: 2 => prénom).
	*	- value:	valeur à rechercher dans la colonne choisie.
	*
	* Description: 		
	*		Méthode permettant de rechercher un ou des compétiteurs 
	*		de la base de données
	*/
	
	func (base Bdd) searchCompetiteur(col_num int, value string){
		
		var id_col string
		var searchValue string
		
		searchValue = fmt.Sprint("'%",value,"%'")	//Mise en forme de la valeur recherchée
		id_col, value = col_id2name(col_num, value)	//Transformation (numéro_de_colonne => id_de_la_colonne)
		
		// REQUÊTE
		base.resultat, base.err = base.db.Query(fmt.Sprint("SELECT * FROM competiteurs WHERE ", id_col, " LIKE ", searchValue))
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(base.err)
		}
		defer base.resultat.Close()
		
		var info [10]string	//Tableau contenant les valeurs de chaque colonnes
		
		//RESULTATS
		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
			
		//AFFICHAGE
		fmt.Println(info[0] + "|" + info[1]+ "|" + info[2]+ "|" + info[3] + "|" + info[4]+ "|" + info[5]+ "|" + info[6]+ "|" + info[7]+ "|" + info[8]+ "|" + info[9])
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

	func (base Bdd) addCompetiteur(comp *Competiteur){
		var test bool 
		// Vérification du format des valeurs entrées.
		test = comp.check()
		
		//Si les valeurs sont bonnes
		if (test) {
			//Ajout du compétiteur
			_, base.err = base.db.Exec("INSERT INTO competiteurs (prenom, nom, sexe, num_license, equipe, epreuve1, annonce1, epreuve2, annonce2) VALUES('" +
			comp.prenom + "','" +
			comp.nom + "','" +
			comp.sexe + "','" +
			comp.num_license + "','" +
			comp.equipe + "','" +
			comp.epreuve1 + "'," +
			strconv.Itoa(comp.annonce1) + ",'" +
			comp.epreuve2 + "'," +
			strconv.Itoa(comp.annonce2) + ")")
		} else {
			log.Fatal(fmt.Sprint("Erreur lors de l'ajout du compétiteur ",comp.prenom," ",comp.nom,". données entrées eronnées."))
		}
		
		if base.err != nil {
			fmt.Println("Echec lors de l'ajout: \n", base.err)
			} else {
			fmt.Println("Ajout validé du compétiteur " + comp.nom +" "+ comp.prenom)
		}
	}
	
	/*
	* 		Bdd.deleteCompetiteur:
	* Paramètres:
	*	- col_num: 	numéro de la colonne sur laquelle effectuer la recherche (1 => id/ 2 => équipe).
	*	- value:	valeur à rechercher dans la colonne choisie.
	*
	* Description: 		
	*		Méthode permettant de supprimer les compétiteurs en fonction des critères
	*		en entrée.
	*/

	func (base Bdd) deleteCompetiteur(col_num int, value string){
		var id_col string
		value = fmt.Sprint("'",value,"'")
		
		// Numéro de colonne => Id colonne
		if col_num==1 {
			id_col = "id"		
		} else if col_num==2{
			id_col = "equipe"		
		}
		
		//Si le numéro de colonne est bon
		if !(col_num < 1 && col_num > 2){
		
			//SUPRESSION DES COMPETITEURS
			_, base.err = base.db.Exec("DELETE FROM competiteurs WHERE " + id_col + " = " + value)
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
	* 		Bdd.resetCompetiteurs:
	* Description: 		
	*		Méthode permettant de supprimer tous les compétiteurs contenus dans la base de
	*		données.
	*/
	
	func (base Bdd) resetCompetiteurs(){
		_, base.err = base.db.Exec("DELETE FROM competiteurs")
		if base.err != nil {
			fmt.Println("Echec lors de la remise à 0 de la table competiteurs. \n", base.err)
		} else {
			_, base.err = base.db.Exec("DELETE FROM sqlite_sequence WHERE name='competiteurs'")
			if base.err != nil {
				fmt.Println("Echec lors de la remise à 0 de la table competiteurs: \n", base.err)
				} 
			}
		}
	
		
	/*
	* 		Bdd.exportCompetiteur:
	* Description: 		
	*		Méthode permettant d'exporter un fichier CSV contenant tous les
	*		compétiteurs de la base de données.
	*/
	func (base Bdd) exportCompetiteur(){
	
		base.resultat, base.err = base.db.Query("SELECT * FROM competiteurs")
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
		}
		defer base.resultat.Close()
		
		t := time.Now()
		date := fmt.Sprint(t.Year(),"_",int(t.Month()),"_", t.Day(),"_",t.Hour(),"_", t.Minute(),"_", t.Second())
		
		// Création des fichiers d'archive et d'exploitation.
		file, err := os.Create(fmt.Sprint("export/archives/",date,"-competiteurs.csv"))
		file2, err := os.Create(fmt.Sprint("export/competiteurs.csv"))
		
		if err != nil {
			fmt.Println("Erreur lors de la création du fichier. Avez vous créé un dossier \"export\" dans le dossier de l'application?")
			log.Fatal(err)
		}
		
		var info [10]string
		
		//Ecriture de l'entête (avec \xEF\xBB\xBF pour passer de l'UTF-8 SANS BOM à l'UTF-8)
		file2.WriteString(fmt.Sprint("\xEF\xBB\xBFId; Prenom; Nom; Sexe; Num_License; Equipe; Epreuve1; annonce1; Epreuve2; annonce2\r\n"))
		file.WriteString(fmt.Sprint("\xEF\xBB\xBFId; Prenom; Nom; Sexe; Num_License; Equipe; Epreuve1; annonce1; Epreuve2; annonce2\r\n"))
		
		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
			file.WriteString(fmt.Sprint(info[0],";",info[1],";", info[2],";", info[3],";", info[4],";", info[5],";", info[6],";", info[7],";", info[8],";", info[9],"\r\n"))
			file2.WriteString(fmt.Sprint(info[0],";",info[1],";", info[2],";", info[3],";", info[4],";", info[5],";", info[6],";", info[7],";", info[8],";", info[9],"\r\n"))
		}
	}
	
	
	/*
	* 		Bdd.importCompetiteur:
	* Paramètres:
	*	- chemin: 	Chemin du fichier à importer avec le nom du fichier et l'extension.
	*
	* Description: 		
	*		Méthode permettant d'importer les compétiteurs contenu dans le fichier CSV
	*		contenu dan: "import/competiteurs.csv"
	*/
	
	func (base Bdd) importCompetiteur(){
		// OUVERTURE DU FICHIER
		file, err := os.Open("import/competiteurs.csv")
		if err != nil {
			fmt.Println("Impossible d'ouvrir le fichier \"competiteurs.csv\" dans le dossier import")
			log.Fatal(err)
		}
		defer file.Close()
		var firstCall bool
		
		firstCall = true
		
		//SCAN DU FICHIER
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			info := strings.Split(scanner.Text(), ";")
			if !firstCall{	//FIRSTCALL => PREMIERE LIGNE => EN-TÊTE!
				temps1,errr := strconv.Atoi(info[6])
				temps2,er := strconv.Atoi(info[8])
				if er != nil {
				log.Fatal(er)
				}
				if errr != nil {
				log.Fatal(errr)
				}
				comp := newCompetiteur(0, info[0], info[1], info[2], info[3], info[4], info[5], temps1, info[7],temps2)
				base.addCompetiteur(comp)
			}
			firstCall = false
		}
			

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// Verification de l'unicité de chaque compétiteurs
		base.uniqueness()	
		//Importation equipe dans la table de classement par équipe
		base.importEquipe()
	}
	
	/*
	* 		Bdd.modifCompetiteur:
	* Paramètres:
	*	- id_comp:	id du compétiteur à modifier
	*	- col_num:  Numéro de la colonne sur laquelle effectuer la modification (ex: 2 => prénom).
	*	- newvalue:	Nouvelle valeur à entrée pour la colonne choisie.
	*
	* Description: 		
	*		Méthode permettant de modifier une valeur d'un compétiteur de la base de données.
	*/
	
	func (base Bdd) modifCompetiteur (id_comp int, col_num int, newvalue string){
		var test = true
		// Colonne num => colonne_ID
		col_id, value := col_id2name(col_num, newvalue)
		
		// Vérification du format de la nouvelle valeur:
		test = verifValue(col_num, newvalue)
		
		//Si la valeur est bonne:
		if (test){
			_, base.err = base.db.Exec("UPDATE competiteurs SET " + col_id + " = " + value + " WHERE id = " + strconv.Itoa(id_comp))
			
			if base.err != nil {
				fmt.Println("Echec lors de la modification: \n", base.err)
			} else {
				fmt.Println("Modification du competiteur " + strconv.Itoa(id_comp) + " avec " + col_id + " = " + value)
			}
		} else {
			fmt.Println("Erreur lors de la modifications du compétiteur!")
		}
	}
	
	
	/*
	* 		Bdd.verifValue:
	* Paramètres:
	*	- col_num:  Numéro de la colonne pour laquelle on vérifie la valeur.
	*	- value:	Valeur à vérifier.
	*
	* Description: 		
	*		Méthode permettant de vérifier le format d'une valeur en fonction
	*		de la colonne choisie.
	*/
	func verifValue(col_num int, value string)(bool){
		var verif = true
		verif = true
		switch col_num{
		    case 2, 3:
				match, _ := regexp.MatchString("^[\\p{L}- ]*$", value )
				if(!match){
					verif =false
					fmt.Println("Erreur! Format du prénom.")
				}
			case 4:
				match, _ := regexp.MatchString("([F|H])", value )
				if(!match || len(value) > 1){
					verif =false
					fmt.Println("Erreur! Format du sexe.")
				}
			case 5:
				match, _ := regexp.MatchString("^[A-Za-z0-9-]*$", value )
				if(!match){
					verif =false
					fmt.Println("Erreur! Format du numéro de license.")
				}
			case 6:
				match, _ := regexp.MatchString("^[\\p{L}0-9- _]*$", value )
				if(!match){
					verif =false
					fmt.Println("Erreur! Format du nom d'équipe.")
				}
			case 7,9:
				if(value!="sta" && value!="spd" && value!="dwf" && value!="dnf" && value!="850"){
					verif =false
					fmt.Println("Erreur! Format du epreuve (Rappel des valeurs possibles: sta, spd, dwf, dnf, 850).")
				}
			case 8,10:
				match, _ := regexp.MatchString("(^[0-9]*$)", value)
				if(!match){
					verif = false
					fmt.Println("Erreur! Format du format de l'annonce.")
				}
			default:
				log.Fatal("Numéro de colone invalide")
			}
		return verif
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
	func col_id2name(col_num int, value string)(string, string){
		var col_id string
		
		switch col_num{
		    case 1:
				col_id = "id"
				value = fmt.Sprint("'",value,"'")
			case 2:
				col_id = "prenom"
				value = fmt.Sprint("'",value,"'")
			case 3:
				col_id = "nom"
				value = fmt.Sprint("'",value,"'")
			case 4:
				col_id = "sexe"
				value = fmt.Sprint("'",value,"'")
			case 5:
				col_id = "num_license"
				value = fmt.Sprint("'",value,"'")
			case 6:
				col_id = "equipe"
				value = fmt.Sprint("'",value,"'")
			case 7:
				col_id = "epreuve1"
				value = fmt.Sprint("'",value,"'")
			case 8:
				col_id = "annonce1"
				value = fmt.Sprint("'",value,"'")
			case 9:
				col_id = "epreuve2"
				value = fmt.Sprint("'",value,"'")
			case 10:
				col_id = "annonce2"
				value = fmt.Sprint("'",value,"'")
			default:
				log.Fatal("Numéro invalide")
			}
		return col_id, value
	}
	
			
	/*
	* 		Bdd.uniqueness:
	* Description: Méthode permettent de vérifier l'unicité des champs id et licence censé être unique	
	*		
	*/
	func (base Bdd) uniqueness(){	
	
		base.resultat, base.err = base.db.Query(fmt.Sprint("SELECT * FROM competiteurs"))
		if base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(base.err)
		}
		//defer base.resultat.Close()
		
		var info [10]string

		for base.resultat.Next() {
			base.err = base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9])
			if base.err != nil {
				fmt.Println("Erreur lors de la récupération des résultats: \n")
				log.Fatal(base.err)
			}
			base.verif(info[0],1)
			base.verif(info[4],2)
		}	
	}