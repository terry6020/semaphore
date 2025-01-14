package projects

import (
	"errors"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/lib"
	"github.com/ansible-semaphore/semaphore/util"
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

type publicAlias struct {
	URL string `json:"url"`
}

func getPublicAlias(alias string) publicAlias {

	if alias == "" {
		return publicAlias{
			URL: "",
		}
	}

	aliasURL := util.Config.WebHost

	if !strings.HasSuffix(aliasURL, "/") {
		aliasURL += "/"
	}

	aliasURL += "api/integrations/" + alias

	return publicAlias{
		URL: aliasURL,
	}
}

func GetIntegrationAlias(w http.ResponseWriter, r *http.Request) {
	integration := context.Get(r, "integration").(db.Integration)

	alias, err := helpers.Store(r).GetIntegrationAlias(integration.ProjectID, &integration.ID)

	if err != nil && !errors.Is(err, db.ErrNotFound) {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, getPublicAlias(alias.Alias))
}

func AddIntegrationAlias(w http.ResponseWriter, r *http.Request) {
	integration := context.Get(r, "integration").(db.Integration)
	alias, err := helpers.Store(r).CreateIntegrationAlias(db.IntegrationAlias{
		Alias:         lib.RandomString(16),
		ProjectID:     integration.ProjectID,
		IntegrationID: &integration.ID,
	})

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, getPublicAlias(alias.Alias))
}

func UpdateIntegrationAlias(w http.ResponseWriter, r *http.Request) {
	integration := context.Get(r, "integration").(db.Integration)

	err := helpers.Store(r).UpdateIntegrationAlias(db.IntegrationAlias{
		Alias:         lib.RandomString(16),
		ProjectID:     integration.ProjectID,
		IntegrationID: &integration.ID,
	})

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	alias, err := helpers.Store(r).GetIntegrationAlias(integration.ProjectID, &integration.ID)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, getPublicAlias(alias.Alias))
}

func RemoveIntegrationAlias(w http.ResponseWriter, r *http.Request) {
	integration := context.Get(r, "integration").(db.Integration)

	err := helpers.Store(r).DeleteIntegrationAlias(integration.ProjectID, &integration.ID)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
