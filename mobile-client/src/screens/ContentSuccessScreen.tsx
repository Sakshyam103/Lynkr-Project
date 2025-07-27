import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';

export default function ContentResultScreen({ route, navigation }: any) {
    const { success, message, eventId } = route.params;

    return (
        <View style={styles.container}>
            <View style={styles.resultCard}>
                <Text style={styles.resultIcon}>
                    {success ? '✅' : '❌'}
                </Text>

                <Text style={[styles.resultTitle, success ? styles.successTitle : styles.errorTitle]}>
                    {success ? 'Content Created Successfully!' : 'Upload Failed'}
                </Text>

                <Text style={styles.resultMessage}>
                    {message || (success ? 'Your content has been uploaded and is now live.' : 'Failed to upload content. Please try again.')}
                </Text>

                <TouchableOpacity
                    style={[styles.continueButton, success ? styles.successButton : styles.errorButton]}
                    onPress={() => navigation.navigate('EventDetails', { eventId })}
                >
                    <Text style={styles.continueText}>
                        {success ? 'Back to Event' : 'Try Again'}
                    </Text>
                </TouchableOpacity>
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#f8f9fa',
        justifyContent: 'center',
        alignItems: 'center',
        padding: 20,
    },
    resultCard: {
        backgroundColor: '#fff',
        borderRadius: 20,
        padding: 40,
        alignItems: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.1,
        shadowRadius: 12,
        elevation: 4,
        width: '100%',
        maxWidth: 350,
    },
    resultIcon: {
        fontSize: 80,
        marginBottom: 20,
    },
    resultTitle: {
        fontSize: 24,
        fontWeight: 'bold',
        textAlign: 'center',
        marginBottom: 16,
    },
    successTitle: {
        color: '#28a745',
    },
    errorTitle: {
        color: '#dc3545',
    },
    resultMessage: {
        fontSize: 16,
        color: '#666',
        textAlign: 'center',
        lineHeight: 22,
        marginBottom: 30,
    },
    continueButton: {
        paddingHorizontal: 40,
        paddingVertical: 16,
        borderRadius: 12,
        width: '100%',
    },
    successButton: {
        backgroundColor: '#28a745',
    },
    errorButton: {
        backgroundColor: '#dc3545',
    },
    continueText: {
        color: '#fff',
        fontSize: 16,
        fontWeight: '600',
        textAlign: 'center',
    },
});
