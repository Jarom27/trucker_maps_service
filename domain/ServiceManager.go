package domain

import (
	"sync"
	"time"
)

// GPSData representa la última posición de un dispositivo.
type GPSData struct {
	DeviceID  string
	Lat       float64
	Lng       float64
	Timestamp time.Time
}

// GPSManager maneja el estado y la lógica de negocio de los dispositivos GPS.
type GPSManager struct {
	data map[string]GPSData // Almacena la última posición de cada dispositivo
	mu   sync.Mutex         // Mutex para proteger el acceso concurrente
}

// NewGPSManager crea una nueva instancia de GPSManager.
func NewGPSManager() *GPSManager {
	return &GPSManager{
		data: make(map[string]GPSData),
	}
}

// UpdateGPSData actualiza la última posición de un dispositivo.
func (gm *GPSManager) UpdateGPSData(deviceID string, lat, lng float64) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	gm.data[deviceID] = GPSData{
		DeviceID:  deviceID,
		Lat:       lat,
		Lng:       lng,
		Timestamp: time.Now(),
	}
}

// GetAllGPSData devuelve todos los datos actuales.
func (gm *GPSManager) GetAllGPSData() map[string]GPSData {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Crear una copia para evitar la exposición del mapa original.
	copiedData := make(map[string]GPSData)
	for k, v := range gm.data {
		copiedData[k] = v
	}

	return copiedData
}

// WorkerPool para procesar mensajes de RabbitMQ
func (gm *GPSManager) StartWorkerPool(workerCount int, jobs <-chan GPSData) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for job := range jobs {
				gm.UpdateGPSData(job.DeviceID, job.Lat, job.Lng)
				// Aquí puedes agregar lógica adicional, como validar datos
			}
		}(i)
	}
}
