package domain

import (
	"sync"
	"time"
	"trucker_maps_service/models"
)

// GPSManager maneja el estado y la lógica de negocio de los dispositivos GPS.
type GPSManager struct {
	data map[string]models.GPSData // Almacena la última posición de cada dispositivo
	mu   sync.Mutex                // Mutex para proteger el acceso concurrente
}

// NewGPSManager crea una nueva instancia de GPSManager.
func NewGPSManager() *GPSManager {
	return &GPSManager{
		data: make(map[string]models.GPSData),
	}
}

// UpdateGPSData actualiza la última posición de un dispositivo.
func (gm *GPSManager) UpdateGPSData(deviceID string, lat, lng float64) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	gm.data[deviceID] = models.GPSData{
		Device_id: deviceID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now(),
	}
}

// GetAllGPSData devuelve todos los datos actuales.
func (gm *GPSManager) GetAllGPSData() map[string]models.GPSData {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Crear una copia para evitar la exposición del mapa original.
	copiedData := make(map[string]models.GPSData)
	for k, v := range gm.data {
		copiedData[k] = v
	}

	return copiedData
}

// WorkerPool para procesar mensajes de RabbitMQ
func (gm *GPSManager) StartWorkerPool(workerCount int, jobs <-chan models.GPSData) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for job := range jobs {
				gm.UpdateGPSData(job.Device_id, job.Latitude, job.Longitude)
				// Aquí puedes agregar lógica adicional, como validar datos
			}
		}(i)
	}
}
