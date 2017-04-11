package main

import (
	"strconv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"bufio"
	"os"
	"strings"
	"sort"
	)


type Planning struct
{
	base *Bdd
	comp []*Competiteur
	cfgEpreuves []*ConfigurationEpreuve
	planEpreuves []*PlanningEpreuve
	
}

func (p *Planning) getCompetiteur(){
	p.base.resultat, p.base.err = p.base.db.Query("SELECT * FROM competiteurs")
		if p.base.err != nil {
			fmt.Println("Erreur lors de l'execution de la requête")
			log.Fatal(p.base.err)
		}
		defer p.base.resultat.Close()
		
		var info [10]string
		var temps1 int
		var temps2 int
		
		var nextcomp *Competiteur 
		
		for p.base.resultat.Next() {
			p.base.err = p.base.resultat.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5], &info[6], &info[7], &info[8], &info[9])
			if p.base.err != nil {
			fmt.Println("Erreur lors de la récupération des résultats: \n")
			log.Fatal(p.base.err)
			}
		temps1,_ = strconv.Atoi(info[7])
		temps2,_ = strconv.Atoi(info[9])
		
		nextcomp = newCompetiteur(info[0], info[1], info[2], info[3], info[4], info[5], info[6], temps1, info[8], temps2)
		p.comp = append(p.comp, nextcomp)
		}
}


// Enregistrer les épreuves dans le tableau
func (p *Planning) getConfigurationEpreuve(){
	file, err := os.Open("Configuration/Configuration.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	
	var firstCall bool
	firstCall = true
	var nextconfig *ConfigurationEpreuve
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), ";")
		if !firstCall{
		ordre, _ := strconv.Atoi(info[0])
		seuilMin, _ := strconv.Atoi(info[2])
		seuilMax, _ := strconv.Atoi(info[3])
		nbPassages, _ := strconv.Atoi(info[4])
		marge, _ := strconv.Atoi(info[5])

		
		nextconfig = newConfigurationEpreuve(ordre, info[1], seuilMin, seuilMax,nbPassages, marge, info[6])
		p.cfgEpreuves = append(p.cfgEpreuves, nextconfig)
		}
		firstCall = false
	}
	sort.Sort(triEpreuves(p.cfgEpreuves))
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}



func (p *Planning) EpGeneration(numEp int){
	numEp = numEp - 1
	if numEp < len(p.cfgEpreuves){
		for j := 0; j < len(p.comp); j++ {
			if p.comp[j].epreuve1 == p.cfgEpreuves[numEp].id{
				plannEp := newPlanningEpreuve(p.cfgEpreuves[numEp].id, p.comp[j].id, p.comp[j].prenom, p.comp[j].nom, p.comp[j].temps1, "00h00")
				p.planEpreuves = append(p.planEpreuves,plannEp)

			} else if p.comp[j].epreuve2 == p.cfgEpreuves[numEp].id{
				plannEp := newPlanningEpreuve(p.cfgEpreuves[numEp].id, p.comp[j].id, p.comp[j].prenom, p.comp[j].nom, p.comp[j].temps2, "00h00")
				p.planEpreuves = append(p.planEpreuves,plannEp)
			}
		}
		sort.Sort(triAnnonces(p.planEpreuves))
	} else{
		fmt.Println("Numéro d'épreuve invalide")
	}
}


func (p *Planning) generationHoraires(numEp int){
	var nbCompPassage int = 0
	nbCompPassage = 0
	var annonceMax int
	annonceMax = 0


	numEp = numEp - 1
	if numEp < len(p.cfgEpreuves){
		heureSTR := strings.Split(p.cfgEpreuves[numEp].heureDebut, "h")
		var heure []int
		h,_ := strconv.Atoi(heureSTR[0])
		m,_ := strconv.Atoi(heureSTR[1])
		heure = append(heure, h, m)
		
		var UMesure string
		UMesure = "m"
		
		if p.cfgEpreuves[numEp].id!="DNF" && p.cfgEpreuves[numEp].id!="DWF"{
			UMesure = "s"
		}
		
		
		for j := 0; j < len(p.planEpreuves); j++ {
			if p.planEpreuves[j].idEpreuve == p.cfgEpreuves[numEp].id{
				//1 er competiteurs => Plus grosse annonce
				if nbCompPassage == 0 && UMesure == "s"{
					//Passage sec => Min
					annonceMax = p.planEpreuves[j].annonce/60
					if (p.planEpreuves[j].annonce%60) != 0{
						annonceMax = annonceMax + 1
					}
				}
				//Configuration de l'heure
				if (heure[1] < 10){
					p.planEpreuves[j].heurePassage = fmt.Sprint(heure[0],"h0",heure[1])
				} else{
					p.planEpreuves[j].heurePassage = fmt.Sprint(heure[0],"h",heure[1])
				}
				
				nbCompPassage = nbCompPassage + 1
				
				//Si le nombre de compétiteur est plein pour ce passage:
				// On modifie l'heure du prochain passage (avec l'annonce max et la marge)
				
				if nbCompPassage == p.cfgEpreuves[numEp].nbPassages{
					nbCompPassage = 0
					heure[1] = heure[1] + annonceMax + p.cfgEpreuves[numEp].marge
					for heure[1] >= 60{
						heure[1] = heure[1] - 60
						heure[0] = heure[0]	+ 1				
					}
				
				}
			}
		}
	} else{
		fmt.Println("Numéro d'épreuve invalide")
	}
}

func (p *Planning) generationPlanning(){
	file, err := os.Create(fmt.Sprint("planning.csv"))
			if err != nil {
				fmt.Println("Erreur lors de la création du fichier planning:\n")
				log.Fatal(err)
			}
	file.WriteString(fmt.Sprint("Epreuve; Id Competiteur; Prenom; Nom; Annonce(s/m); Heure de passage\r\n"))
			
	for j := 1; j <= len(p.cfgEpreuves); j++ {
		p.planEpreuves = p.planEpreuves[:0]
		p.EpGeneration(j)
		p.generationHoraires(j)
		p.displayPlanningEpreuve()
		p.exportPlanEpreuve()
	}

}


func (p *Planning) exportPlanEpreuve(){
	file, err := os.OpenFile(fmt.Sprint("planning.csv"),os.O_APPEND|os.O_WRONLY, 0777)
			if err != nil {
				fmt.Println("Erreur lors de la création du fichier")
				log.Fatal(err)
			}

			fmt.Println(len(p.planEpreuves))
			for j := 0; j < len(p.planEpreuves); j++ {
						file.WriteString(fmt.Sprint(p.planEpreuves[j].idEpreuve,";",p.planEpreuves[j].idComp,";",p.planEpreuves[j].prenom,";",p.planEpreuves[j].nom,";",strconv.Itoa(p.planEpreuves[j].annonce),";",p.planEpreuves[j].heurePassage,"\r\n"))
			}

}

func (p Planning) displayCompetiteur(){
	for j := 0; j < len(p.comp); j++ {
		p.comp[j].display()
    }
}

func (p Planning) displayConfigurationEpreuve(){
	for j := 0; j < len(p.cfgEpreuves); j++ {
		p.cfgEpreuves[j].display()
    }
}

func (p Planning) displayPlanningEpreuve(){
	for j := 0; j < len(p.planEpreuves); j++ {
		p.planEpreuves[j].display()
    }
}

func newPlanning(cheminBdd string)(*Planning){
	p := new(Planning)
	p.base = newBdd(cheminBdd)	
	
	return p
}