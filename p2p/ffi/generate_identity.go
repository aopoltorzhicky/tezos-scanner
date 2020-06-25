package ffi

/*
#cgo LDFLAGS: -L. ${SRCDIR}/libtezos.so -ldl
#include <stdio.h>
#include <stdlib.h>
#include <caml/alloc.h>
#include <caml/callback.h>
#include <caml/memory.h>
#include <caml/mlvalues.h>
value ml_context_dir_mem(value context_hash, value block_hash, value operation_hash, value key, value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, time_period);
}
value ml_context_remove_rec(value context_hash, value block_hash, value operation_hash, value key, value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, time_period);
}
value ml_context_commit(value parent_context_hash, value block_hash, value new_context_hash, value time_period) {
 CAMLparam4(parent_context_hash, block_hash, new_context_hash, time_period);
}
value ml_context_raw_get(value context_hash, value block_hash, value operation_hash, value key, value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, time_period);
}
value ml_context_delete(value context_hash, value block_hash, value operation_hash, value key, value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, time_period);
}
value ml_context_set(value context_hash, value block_hash, value operation_hash, value key, value value_and_json) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, value_and_json);
}
value ml_context_mem(value context_hash, value block_hash, value operation_hash, value key, value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, time_period);
}
value ml_context_copy(value context_hash, value block_hash, value operation_hash, value from_to_key,
                     value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, from_to_key, time_period);
}
value ml_context_fold(value context_hash, value block_hash, value operation_hash, value key, value time_period) {
 CAMLparam5(context_hash, block_hash, operation_hash, key, time_period);
}
value ml_context_checkout(value context_hash, value time_period) { CAMLparam2(context_hash, time_period); }

const char* valueToString(value val) {
    return String_val(val);
}

// https://dev.to/mattn/call-go-function-from-c-function-1n3
*/
import "C"
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"unsafe"
)

const identityFile = "identity.json"

// Identity -
type Identity struct {
	PeerID           string `json:"peer_id"`
	PublicKey        string `json:"public_key"`
	SecretKey        string `json:"secret_key"`
	ProofOfWorkStamp string `json:"proof_of_work_stamp"`
}

// GetIdentity -
func GetIdentity() (Identity, error) {
	_, err := os.Stat(identityFile)
	if err != nil {
		if os.IsNotExist(err) {
			identity, err := GenerateIdentity()
			if err != nil {
				log.Println(err)
			}
			log.Println("Identity Generation Successful")
			log.Println("peerId: ", identity.PeerID)
			log.Println("publicKey: ", identity.PublicKey)
			log.Println("secretKey: ", identity.SecretKey)
			log.Println("POW: ", identity.ProofOfWorkStamp)
			return identity, nil
		}
		return Identity{}, err
	}
	file, _ := ioutil.ReadFile(identityFile)
	identity := Identity{}
	err = json.Unmarshal([]byte(file), &identity)
	if err != nil {
		return Identity{}, err
	}
	log.Println("Identity obtained successfully")
	log.Println("peerId: ", identity.PeerID)
	log.Println("publicKey: ", identity.PublicKey)
	log.Println("secretKey: ", identity.SecretKey)
	log.Println("POW: ", identity.ProofOfWorkStamp)
	return identity, nil
}

// GenerateIdentity -
func GenerateIdentity() (Identity, error) {
	log.Println("Generate Identity...")
	path := C.CString("./libtezos.so")
	defer C.free(unsafe.Pointer(path))
	var argv **C.char
	defer C.free(unsafe.Pointer(argv))
	argv = &path
	C.caml_startup(argv)
	funcName := C.CString("generate_identity")
	defer C.free(unsafe.Pointer(funcName))
	funcPoint := C.caml_named_value(funcName)
	callbackFunc := C.caml_callback(*funcPoint, C.copy_double(26.))
	returnStr := C.valueToString(callbackFunc)
	str := C.GoString(returnStr)
	identity := Identity{}
	err := json.Unmarshal([]byte(str), &identity)
	if err != nil {
		return Identity{}, err
	}
	jsonData, err := json.Marshal(identity)
	if err != nil {
		return Identity{}, err
	}
	jsonFile, err := os.Create(identityFile)
	if err != nil {
		return Identity{}, err
	}
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return Identity{}, err
	}
	return identity, nil
}
