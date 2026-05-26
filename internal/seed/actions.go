package seed

import (
	"context"
	"log"
	"time"

	"nw-back/internal/postgres/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

func createCustomerActions(ctx context.Context, queries *db.Queries, customers []db.Customer) error {
	actionsCreated := 0

	for _, customer := range customers {
		actionsCount := gofakeit.Number(2, 10)

		for range actionsCount {
			_, err := queries.CreateCustomerAction(ctx, randomCustomerAction(customer))
			if err != nil {
				return err
			}

			actionsCreated++
		}
	}

	log.Printf("%d customer actions created", actionsCreated)
	return nil
}

func randomCustomerAction(customer db.Customer) db.CreateCustomerActionParams {
	actionType := randomCustomerActionType()

	return db.CreateCustomerActionParams{
		CustomerID: customer.ID,
		Type:       actionType,
		Comments:   randomCustomerActionComment(actionType),
		InformantName: pgtype.Text{
			String: gofakeit.Name(),
			Valid:  true,
		},
		ActionDate: pgtype.Timestamptz{
			Time:  randomActionDate(),
			Valid: true,
		},
	}
}

func randomCustomerActionType() db.CustomerActionType {
	actionTypes := []db.CustomerActionType{
		db.CustomerActionTypeCall,
		db.CustomerActionTypeEmail,
		db.CustomerActionTypePersonalVisit,
		db.CustomerActionTypeOther,
	}

	return actionTypes[gofakeit.Number(0, len(actionTypes)-1)]
}

func randomCustomerActionComment(actionType db.CustomerActionType) string {
	comments := map[db.CustomerActionType][]string{
		db.CustomerActionTypeCall: {
			"Se llamo al cliente para revisar el estado de la cuenta.",
			"Contacto telefonico exitoso, solicito reenviar detalle pendiente.",
			"No atendio la llamada, se intentara nuevamente en los proximos dias.",
			"Se converso con administracion sobre los vencimientos abiertos.",
			"El cliente pidio coordinar un nuevo llamado con el area financiera.",
			"Se confirmo por telefono la recepcion de la documentacion enviada.",
			"El contacto indico que revisara internamente el estado de pago.",
			"Se dejo aviso telefonico para seguimiento de la cuenta.",
			"El cliente manifesto conformidad con el servicio actual.",
			"Se acordo mantener contacto telefonico para proximas novedades.",
		},
		db.CustomerActionTypeEmail: {
			"Se envio correo con resumen de cuenta y proximos vencimientos.",
			"El cliente respondio solicitando detalle de facturas anteriores.",
			"Se remitio comprobante y documentacion comercial actualizada.",
			"Correo enviado al contacto administrativo para seguimiento.",
			"Se envio recordatorio formal sobre informacion pendiente.",
			"El cliente confirmo recepcion del correo enviado.",
			"Se compartio detalle mensual para facilitar conciliacion interna.",
			"Se envio propuesta de regularizacion por correo electronico.",
			"Correo derivado al area contable segun solicitud del cliente.",
			"Se dejo constancia escrita de los temas conversados.",
		},
		db.CustomerActionTypePersonalVisit: {
			"Visita realizada para revisar necesidades operativas del cliente.",
			"Se relevaron comentarios del equipo durante la visita presencial.",
			"El cliente solicito una nueva visita para analizar mejoras.",
			"Se presentaron novedades del servicio en reunion presencial.",
			"Visita coordinada con administracion y responsable de pagos.",
			"Se recopilaron observaciones sobre uso y satisfaccion del servicio.",
			"El contacto pidio seguimiento posterior por correo.",
			"Se reviso documentacion pendiente durante la visita.",
			"Reunion presencial productiva con buena predisposicion del cliente.",
			"Se acordo una nueva instancia de seguimiento comercial.",
		},
		db.CustomerActionTypeOther: {
			"Se actualizo informacion interna relevante para la cuenta.",
			"Cliente marcado para seguimiento preventivo durante la semana.",
			"Se reviso historial y se definieron proximos pasos.",
			"Se registro observacion administrativa sin contacto directo.",
			"Cuenta revisada por posibles cambios en el comportamiento de pago.",
			"Se ajusto prioridad de seguimiento segun estado actual.",
			"Se agrego nota interna para el equipo comercial.",
			"Se valido informacion disponible antes de proxima comunicacion.",
			"Se dejo registro de gestion pendiente de confirmacion.",
			"Cuenta incluida en revision operativa periodica.",
		},
	}

	actionComments := comments[actionType]
	return actionComments[gofakeit.Number(0, len(actionComments)-1)]
}

func randomActionDate() time.Time {
	now := time.Now().UTC()
	secondsInLastThreeMonths := int(now.Sub(now.AddDate(0, -3, 0)).Seconds())
	randomSeconds := gofakeit.Number(0, secondsInLastThreeMonths)

	return now.Add(-time.Duration(randomSeconds) * time.Second)
}
